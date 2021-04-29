package identity_platform_service_test

import (
	"firebase.google.com/go/v4/auth"
	mocks "github.com/nitric-dev/membrane/mocks/identityplatform"
	identity_platform_plugin "github.com/nitric-dev/membrane/plugins/auth/identityplatform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Identityplatform", func() {
	mockFirebaseAuth := mocks.NewMockFirebaseAuth()
	authPlugin := identity_platform_plugin.NewWithClient(mockFirebaseAuth)

	AfterEach(func() {
		mockFirebaseAuth.Clear()
	})

	Context("Create", func() {
		When("the user does not already exist", func() {
			It("should successfully create the user", func() {
				err := authPlugin.Create("test", "test", "test@test.com", "test")

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Storing the user")
				testUser := mockFirebaseAuth.GetUser("test", "test")

				Expect(testUser).ToNot(BeNil())
				Expect(testUser.Email).To(Equal("test@test.com"))
			})
		})

		When("A user with the same ID already exists", func() {
			BeforeEach(func() {
				mockFirebaseAuth.SetTenants([]*mocks.MockTenant{
					{
						T: &auth.Tenant{
							DisplayName: "test",
							ID:          "test",
						},
						Users: []*auth.UserRecord{{
							TenantID: "test",
							UserInfo: &auth.UserInfo{
								UID:         "test",
								DisplayName: "test2@test.com",
								Email:       "test2@test.com",
							},
						}},
					},
				})
			})

			It("should return an error", func() {
				err := authPlugin.Create("test", "test", "test@test.com", "test")

				Expect(err).Should(HaveOccurred())
			})

		})

		When("A user with the same email already exists", func() {
			BeforeEach(func() {
				mockFirebaseAuth.SetTenants([]*mocks.MockTenant{
					{
						T: &auth.Tenant{
							DisplayName: "test",
							ID:          "test",
						},
						Users: []*auth.UserRecord{{
							TenantID: "test",
							UserInfo: &auth.UserInfo{
								UID:         "test2",
								DisplayName: "test@test.com",
								Email:       "test@test.com",
							},
						}},
					},
				})
			})

			It("should return an error", func() {
				err := authPlugin.Create("test", "test", "test@test.com", "test")

				Expect(err).Should(HaveOccurred())
			})
		})

	})
})
