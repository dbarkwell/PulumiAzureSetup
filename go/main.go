package main

import (
	"fmt"
	"github.com/pulumi/pulumi-azure/sdk/go/azure/ad"
	"github.com/pulumi/pulumi-azure/sdk/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/go/azure/role"
	"github.com/pulumi/pulumi-azuread/sdk/go/azuread"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi/sdk/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		config := config.New(ctx, "")
		subscriptionId := config.Require("subscriptionId")

		azureADGroup, err := azuread.NewGroup(ctx, "AzureUsersGroup", &azuread.GroupArgs{
			Name: pulumi.String("All Azure Users"),
		})
		if err != nil {
			return err
		}

		scope := fmt.Sprintf("/subscriptions/%s", subscriptionId)
		roleDefinition, err := role.NewDefinition(ctx, "Role", &role.DefinitionArgs{
			Name:        pulumi.String("Register Azure resource providers"),
			Description: pulumi.String("Can register Azure resource providers"),
			Permissions: role.DefinitionPermissionArray{
				role.DefinitionPermissionArgs{Actions: pulumi.StringArray{pulumi.String("*/register/action")}},
			},
			AssignableScopes: pulumi.StringArray{pulumi.String(scope)},
			Scope: pulumi.String(scope),
		})
		if err != nil {
			return err
		}

		_, err = role.NewAssignment(ctx, "RoleAssignment", &role.AssignmentArgs{
			PrincipalId:                  azureADGroup.ObjectId,
			RoleDefinitionName:           roleDefinition.Name,
			Scope:                        pulumi.String(scope),
		})
		if err!= nil {
			return err
		}

		resourceGroupName := "Pulumi"
		_, err = core.NewResourceGroup(ctx, "ResourceGroup", &core.ResourceGroupArgs{
			Location: pulumi.String("CanadaEast"),
			Name:     pulumi.String(resourceGroupName),
		})
		if err != nil {
			return err
		}

		application, err := ad.NewApplication(ctx, "Application", &ad.ApplicationArgs{
			Name: pulumi.String("PulumiAzureSetup"),
		})
		if err != nil {
			return err
		}

		servicePrincipal, err := ad.NewServicePrincipal(ctx, "ServicePrincipal", &ad.ServicePrincipalArgs{ApplicationId: application.ApplicationId})
		if err != nil {
			return err
		}

		_, err = azuread.NewGroupMember(ctx, "GroupMembership", &azuread.GroupMemberArgs{
			GroupObjectId:  azureADGroup.ObjectId,
			MemberObjectId: servicePrincipal.ID(),
		})
		if err != nil {
			return err
		}

		_, err = role.NewAssignment(ctx, "RoleAssignmentRG", &role.AssignmentArgs{
			PrincipalId:        servicePrincipal.ID(),
			RoleDefinitionName: pulumi.String("Contributor"),
			Scope:              pulumi.String(scope),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
