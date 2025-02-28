//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

var (
	err               error
	ctx               context.Context
	cred              azcore.TokenCredential
	location          = getEnv("LOCATION", "westus")
	resourceGroupName = getEnv("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	subscriptionId    = getEnv("AZURE_SUBSCRIPTION_ID", "")
)

func main() {
	ctx = context.Background()
	cred, err = azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	createResourceGroup()
	signalrSample()
	deleteResourceGroup()
}

func signalrSample() {
	var resourceName string
	// From step Generate_Unique_Name
	template := map[string]interface{}{
		"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"outputs": map[string]interface{}{
			"resourceName": map[string]interface{}{
				"type":  "string",
				"value": "[variables('name').value]",
			},
		},
		"resources": []interface{}{},
		"variables": map[string]interface{}{
			"name": map[string]interface{}{
				"type": "string",
				"metadata": map[string]interface{}{
					"description": "Name of the SignalR service.",
				},
				"value": "[concat('sw',uniqueString(resourceGroup().id))]",
			},
		},
	}
	params := map[string]interface{}{}
	deployment := armresources.Deployment{
		Properties: &armresources.DeploymentProperties{
			Template:   template,
			Parameters: params,
			Mode:       to.Ptr(armresources.DeploymentModeIncremental),
		},
	}
	deploymentExtend := createDeployment("Generate_Unique_Name", &deployment)
	resourceName = deploymentExtend.Properties.Outputs.(map[string]interface{})["resourceName"].(map[string]interface{})["value"].(string)

	// From step SignalR_CheckNameAvailability
	signalRClient, err := test.NewSignalRClient(subscriptionId, cred, nil)
	if err != nil {
		panic(err)
	}
	_, err = signalRClient.CheckNameAvailability(ctx,
		location,
		test.NameAvailabilityParameters{
			Name: to.Ptr(resourceName),
			Type: to.Ptr("Microsoft.SignalRService/SignalR"),
		},
		nil)
	if err != nil {
		panic(err)
	}

	// From step SignalR_CreateOrUpdate
	signalRClientCreateOrUpdateResponsePoller, err := signalRClient.BeginCreateOrUpdate(ctx,
		resourceGroupName,
		resourceName,
		test.ResourceInfo{
			Location: to.Ptr(location),
			Tags: map[string]*string{
				"key1": to.Ptr("value1"),
			},
			Identity: &test.ManagedIdentity{
				Type: to.Ptr(test.ManagedIdentityTypeSystemAssigned),
			},
			Kind: to.Ptr(test.ServiceKindSignalR),
			Properties: &test.SignalRProperties{
				Cors: &test.SignalRCorsSettings{
					AllowedOrigins: []*string{
						to.Ptr("https://foo.com"),
						to.Ptr("https://bar.com")},
				},
				DisableAADAuth:   to.Ptr(false),
				DisableLocalAuth: to.Ptr(false),
				Features: []*test.SignalRFeature{
					{
						Flag:       to.Ptr(test.FeatureFlagsServiceMode),
						Properties: map[string]*string{},
						Value:      to.Ptr("Serverless"),
					},
					{
						Flag:       to.Ptr(test.FeatureFlagsEnableConnectivityLogs),
						Properties: map[string]*string{},
						Value:      to.Ptr("True"),
					},
					{
						Flag:       to.Ptr(test.FeatureFlagsEnableMessagingLogs),
						Properties: map[string]*string{},
						Value:      to.Ptr("False"),
					},
					{
						Flag:       to.Ptr(test.FeatureFlagsEnableLiveTrace),
						Properties: map[string]*string{},
						Value:      to.Ptr("False"),
					}},
				NetworkACLs: &test.SignalRNetworkACLs{
					DefaultAction: to.Ptr(test.ACLActionDeny),
					PrivateEndpoints: []*test.PrivateEndpointACL{
						{
							Allow: []*test.SignalRRequestType{
								to.Ptr(test.SignalRRequestTypeServerConnection)},
							Name: to.Ptr(resourceName + ".1fa229cd-bf3f-47f0-8c49-afb36723997e"),
						}},
					PublicNetwork: &test.NetworkACL{
						Allow: []*test.SignalRRequestType{
							to.Ptr(test.SignalRRequestTypeClientConnection)},
					},
				},
				PublicNetworkAccess: to.Ptr("Enabled"),
				TLS: &test.SignalRTLSSettings{
					ClientCertEnabled: to.Ptr(false),
				},
				Upstream: &test.ServerlessUpstreamSettings{
					Templates: []*test.UpstreamTemplate{
						{
							Auth: &test.UpstreamAuthSettings{
								Type: to.Ptr(test.UpstreamAuthTypeManagedIdentity),
								ManagedIdentity: &test.ManagedIdentitySettings{
									Resource: to.Ptr("api://example"),
								},
							},
							CategoryPattern: to.Ptr("*"),
							EventPattern:    to.Ptr("connect,disconnect"),
							HubPattern:      to.Ptr("*"),
							URLTemplate:     to.Ptr("https://example.com/chat/api/connect"),
						}},
				},
			},
			SKU: &test.ResourceSKU{
				Name:     to.Ptr("Standard_S1"),
				Capacity: to.Ptr[int32](1),
				Tier:     to.Ptr(test.SignalRSKUTierStandard),
			},
		},
		nil)
	if err != nil {
		panic(err)
	}
	_, err = signalRClientCreateOrUpdateResponsePoller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}

	// From step SignalR_Get
	_, err = signalRClient.Get(ctx,
		resourceGroupName,
		resourceName,
		nil)
	if err != nil {
		panic(err)
	}

	// From step SignalR_Update
	signalRClientUpdateResponsePoller, err := signalRClient.BeginUpdate(ctx,
		resourceGroupName,
		resourceName,
		test.ResourceInfo{
			Location: to.Ptr(location),
			Tags: map[string]*string{
				"key1": to.Ptr("value1"),
			},
			Identity: &test.ManagedIdentity{
				Type: to.Ptr(test.ManagedIdentityTypeSystemAssigned),
			},
			Kind: to.Ptr(test.ServiceKindSignalR),
			Properties: &test.SignalRProperties{
				Cors: &test.SignalRCorsSettings{
					AllowedOrigins: []*string{
						to.Ptr("https://foo.com"),
						to.Ptr("https://bar.com")},
				},
				DisableAADAuth:   to.Ptr(false),
				DisableLocalAuth: to.Ptr(false),
				Features: []*test.SignalRFeature{
					{
						Flag:       to.Ptr(test.FeatureFlagsServiceMode),
						Properties: map[string]*string{},
						Value:      to.Ptr("Serverless"),
					},
					{
						Flag:       to.Ptr(test.FeatureFlagsEnableConnectivityLogs),
						Properties: map[string]*string{},
						Value:      to.Ptr("True"),
					},
					{
						Flag:       to.Ptr(test.FeatureFlagsEnableMessagingLogs),
						Properties: map[string]*string{},
						Value:      to.Ptr("False"),
					},
					{
						Flag:       to.Ptr(test.FeatureFlagsEnableLiveTrace),
						Properties: map[string]*string{},
						Value:      to.Ptr("False"),
					}},
				NetworkACLs: &test.SignalRNetworkACLs{
					DefaultAction: to.Ptr(test.ACLActionDeny),
					PrivateEndpoints: []*test.PrivateEndpointACL{
						{
							Allow: []*test.SignalRRequestType{
								to.Ptr(test.SignalRRequestTypeServerConnection)},
							Name: to.Ptr(resourceName + ".1fa229cd-bf3f-47f0-8c49-afb36723997e"),
						}},
					PublicNetwork: &test.NetworkACL{
						Allow: []*test.SignalRRequestType{
							to.Ptr(test.SignalRRequestTypeClientConnection)},
					},
				},
				PublicNetworkAccess: to.Ptr("Enabled"),
				TLS: &test.SignalRTLSSettings{
					ClientCertEnabled: to.Ptr(false),
				},
				Upstream: &test.ServerlessUpstreamSettings{
					Templates: []*test.UpstreamTemplate{
						{
							Auth: &test.UpstreamAuthSettings{
								Type: to.Ptr(test.UpstreamAuthTypeManagedIdentity),
								ManagedIdentity: &test.ManagedIdentitySettings{
									Resource: to.Ptr("api://example"),
								},
							},
							CategoryPattern: to.Ptr("*"),
							EventPattern:    to.Ptr("connect,disconnect"),
							HubPattern:      to.Ptr("*"),
							URLTemplate:     to.Ptr("https://example.com/chat/api/connect"),
						}},
				},
			},
			SKU: &test.ResourceSKU{
				Name:     to.Ptr("Standard_S1"),
				Capacity: to.Ptr[int32](1),
				Tier:     to.Ptr(test.SignalRSKUTierStandard),
			},
		},
		nil)
	if err != nil {
		panic(err)
	}
	_, err = signalRClientUpdateResponsePoller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}

	// From step SignalR_ListKeys
	_, err = signalRClient.ListKeys(ctx,
		resourceGroupName,
		resourceName,
		nil)
	if err != nil {
		panic(err)
	}

	// From step SignalR_RegenerateKey
	signalRClientRegenerateKeyResponsePoller, err := signalRClient.BeginRegenerateKey(ctx,
		resourceGroupName,
		resourceName,
		test.RegenerateKeyParameters{
			KeyType: to.Ptr(test.KeyTypePrimary),
		},
		nil)
	if err != nil {
		panic(err)
	}
	_, err = signalRClientRegenerateKeyResponsePoller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}

	// From step SignalR_Restart
	signalRClientRestartResponsePoller, err := signalRClient.BeginRestart(ctx,
		resourceGroupName,
		resourceName,
		nil)
	if err != nil {
		panic(err)
	}
	_, err = signalRClientRestartResponsePoller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}

	// From step Usages_List
	usagesClient, err := test.NewUsagesClient(subscriptionId, cred, nil)
	if err != nil {
		panic(err)
	}
	usagesClientNewListPagerPager := usagesClient.NewListPager(location,
		nil)
	for usagesClientNewListPagerPager.More() {
		_, err := usagesClientNewListPagerPager.NextPage(ctx)
		if err != nil {
			panic(err)
		}
		break
	}

	// From step SignalR_ListByResourceGroup
	signalRClientNewListByResourceGroupPagerPager := signalRClient.NewListByResourceGroupPager(resourceGroupName,
		nil)
	for signalRClientNewListByResourceGroupPagerPager.More() {
		_, err := signalRClientNewListByResourceGroupPagerPager.NextPage(ctx)
		if err != nil {
			panic(err)
		}
		break
	}

	// From step SignalR_ListBySubscription
	signalRClientNewListBySubscriptionPagerPager := signalRClient.NewListBySubscriptionPager(nil)
	for signalRClientNewListBySubscriptionPagerPager.More() {
		_, err := signalRClientNewListBySubscriptionPagerPager.NextPage(ctx)
		if err != nil {
			panic(err)
		}
		break
	}

	// From step Operations_List
	operationsClient, err := test.NewOperationsClient(cred, nil)
	if err != nil {
		panic(err)
	}
	operationsClientNewListPagerPager := operationsClient.NewListPager(nil)
	for operationsClientNewListPagerPager.More() {
		_, err := operationsClientNewListPagerPager.NextPage(ctx)
		if err != nil {
			panic(err)
		}
		break
	}

	// From step SignalR_Delete
	signalRClientDeleteResponsePoller, err := signalRClient.BeginDelete(ctx,
		resourceGroupName,
		resourceName,
		nil)
	if err != nil {
		panic(err)
	}
	_, err = signalRClientDeleteResponsePoller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}
}

func createResourceGroup() error {
	rand.Seed(time.Now().UnixNano())
	resourceGroupName = fmt.Sprintf("go-sdk-sample-%d", rand.Intn(1000))
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		panic(err)
	}
	param := armresources.ResourceGroup{
		Location: to.Ptr(location),
	}
	_, err = rgClient.CreateOrUpdate(ctx, resourceGroupName, param, nil)
	if err != nil {
		panic(err)
	}
	return nil
}

func deleteResourceGroup() error {
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionId, cred, nil)
	if err != nil {
		panic(err)
	}
	pollerResponse, err := rgClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		panic(err)
	}
	_, err = pollerResponse.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func createDeployment(deploymentName string, deployment *armresources.Deployment) *armresources.DeploymentExtended {
	deployClient, err := armresources.NewDeploymentsClient(subscriptionId, cred, nil)
	if err != nil {
		panic(err)
	}
	poller, err := deployClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		deploymentName,
		*deployment,
		&armresources.DeploymentsClientBeginCreateOrUpdateOptions{},
	)
	if err != nil {
		panic(err)
	}
	res, err := poller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		panic(err)
	}
	return &res.DeploymentExtended
}
