package deploy

import (
	"encoding/json"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/batch"
)

// "containerProperties": {
// 	"image": "",
// 	"vcpus": 0,
// 	"memory": 0,
// 	"command": [
// 		""
// 	],
// 	"jobRoleArn": "",
// 	"executionRoleArn": "",
// 	"volumes": [
// 		{
// 			"host": {
// 				"sourcePath": ""
// 			},
// 			"name": "",
// 			"efsVolumeConfiguration": {
// 				"fileSystemId": "",
// 				"rootDirectory": "",
// 				"transitEncryption": "ENABLED",
// 				"transitEncryptionPort": 0,
// 				"authorizationConfig": {
// 					"accessPointId": "",
// 					"iam": "DISABLED"
// 				}
// 			}
// 		}
// 	],
// 	"environment": [
// 		{
// 			"name": "",
// 			"value": ""
// 		}
// 	],
// 	"mountPoints": [
// 		{
// 			"containerPath": "",
// 			"readOnly": true,
// 			"sourceVolume": ""
// 		}
// 	],
// 	"readonlyRootFilesystem": true,
// 	"privileged": true,
// 	"ulimits": [
// 		{
// 			"hardLimit": 0,
// 			"name": "",
// 			"softLimit": 0
// 		}
// 	],
// 	"user": "",
// 	"instanceType": "",
// 	"resourceRequirements": [
// 		{
// 			"value": "",
// 			"type": "MEMORY"
// 		}
// 	],
// 	"linuxParameters": {
// 		"devices": [
// 			{
// 				"hostPath": "",
// 				"containerPath": "",
// 				"permissions": [
// 					"WRITE"
// 				]
// 			}
// 		],
// 		"initProcessEnabled": true,
// 		"sharedMemorySize": 0,
// 		"tmpfs": [
// 			{
// 				"containerPath": "",
// 				"size": 0,
// 				"mountOptions": [
// 					""
// 				]
// 			}
// 		],
// 		"maxSwap": 0,
// 		"swappiness": 0
// 	},
// 	"logConfiguration": {
// 		"logDriver": "syslog",
// 		"options": {
// 			"KeyName": ""
// 		},
// 		"secretOptions": [
// 			{
// 				"name": "",
// 				"valueFrom": ""
// 			}
// 		]
// 	},
// 	"secrets": [
// 		{
// 			"name": "",
// 			"valueFrom": ""
// 		}
// 	],
// 	"networkConfiguration": {
// 		"assignPublicIp": "DISABLED"
// 	},
// 	"fargatePlatformConfiguration": {
// 		"platformVersion": ""
// 	}
// }

type JobDefinitionContainerProperties struct {
	Image            string   `json:"image"`
	Vcpus            int      `json:"vcpus"`
	Memory           int      `json:"memory"`
	Command          []string `json:"command"`
	JobRoleArn       string   `json:"jobRoleArn"`
	ExecutionRoleArn string   `json:"executionRoleArn"`
	Environment      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"environment"`
}

func (p *NitricAwsPulumiProvider) Job(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Job) error {

	// Tag the image

	// Create a new Iam Role for the job

	containerProperties := pulumi.All().ApplyT(func(args []interface{}) (string, error) {
		jobDefinitionContainerProperties := JobDefinitionContainerProperties{
			Image:            "",
			Vcpus:            0,
			Memory:           0,
			Command:          []string{""},
			JobRoleArn:       "",
			ExecutionRoleArn: "",
		}

		containerPropertiesJson, err := json.Marshal(jobDefinitionContainerProperties)
		if err != nil {
			return "", err
		}

		return string(containerPropertiesJson), nil
	}).(pulumi.StringOutput)

	//
	batch.NewJobDefinition(ctx, name, &batch.JobDefinitionArgs{
		ContainerProperties: containerProperties,

		// TODO: Set tags for job definition discovery
		Type: pulumi.String("container"),
		Tags: pulumi.ToStringMap(tags.Tags(p.StackId, name, "job")),
	})

	return nil
}
