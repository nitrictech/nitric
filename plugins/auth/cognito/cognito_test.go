package cognito_service_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	mock "github.com/nitric-dev/membrane/mocks/cognito"
	cognito_plugin "github.com/nitric-dev/membrane/plugins/auth/cognito"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cognito", func() {
	Context("Create", func() {
		mockCognito := mock.NewMockCognitoIdentityProvider()
		authPlugin := cognito_plugin.NewWithClient(mockCognito)

		AfterEach(func() {
			mockCognito.Clear()
		})

		When("the user does not already exist", func() {
			It("should successfully create the user", func() {
				err := authPlugin.Create("test", "test", "test@test.com", "test")

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Creating the missing pool")
				mockPool := mockCognito.GetUserPool("test")
				Expect(mockPool).ToNot(BeNil())

				By("Creating the user")
				user := mockCognito.GetMockUser("test", "test")
				Expect(user).ToNot(BeNil())
				Expect(*user.Username).To(Equal("test"))
				var email *string
				for _, att := range user.Attributes {
					if *att.Name == "email" {
						email = att.Value
					}
				}

				Expect(email).ToNot(BeNil())
				Expect(*email).To(Equal("test@test.com"))
			})
		})

		When("A user with the same id already exists", func() {
			existingUser := &cognitoidentityprovider.UserType{
				Username: aws.String("test"),
			}

			BeforeEach(func() {
				mockCognito.SetPools([]*mock.MockUserPool{
					&mock.MockUserPool{
						Pool: &cognitoidentityprovider.UserPoolType{
							Name: aws.String("test"),
							Id:   aws.String("test"),
						},
						Users: []*cognitoidentityprovider.UserType{
							existingUser,
						},
					},
				})
			})

			It("Should return an error", func() {
				err := authPlugin.Create("test", "test", "test2@test.com", "test")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("username already exists"))
			})
		})

		When("A user with the same email already exists", func() {
			existingUser := &cognitoidentityprovider.UserType{
				Username: aws.String("test"),
				Attributes: []*cognitoidentityprovider.AttributeType{
					&cognitoidentityprovider.AttributeType{
						Name:  aws.String("email"),
						Value: aws.String("test@test.com"),
					},
				},
			}

			BeforeEach(func() {
				mockCognito.SetPools([]*mock.MockUserPool{
					&mock.MockUserPool{
						Pool: &cognitoidentityprovider.UserPoolType{
							Name: aws.String("test"),
							Id:   aws.String("test"),
						},
						Users: []*cognitoidentityprovider.UserType{
							existingUser,
						},
					},
				})
			})

			It("Should return an error", func() {
				err := authPlugin.Create("test", "test2", "test@test.com", "test")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("email"))
			})
		})
	})
})
