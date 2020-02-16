using System.Threading.Tasks;
using Pulumi;
using Pulumi.Azure.Authorization;
using Pulumi.Azure.Authorization.Inputs;
using Pulumi.AzureAD;
using Pulumi.Azure.Core;

class Program
{
    static Task<int> Main()
    {
        return Deployment.RunAsync(() => {
            var config = new Pulumi.Config();
            var subscriptionId = config.Require("subscriptionId");
            
            var azureADGroup = new Group(
                "AzureUsersGroup", 
                new GroupArgs {Name = "All Azure Users"});

            var role = new RoleDefinition(
                "Role",
                new RoleDefinitionArgs
                {
                    Name = "Register Azure resource providers",
                    Description = "Can register Azure resource providers",
                    Permissions = new RoleDefinitionPermissionsArgs {Actions = "*/register/action"},
                    AssignableScopes = $"/subscriptions/{subscriptionId}",
                    Scope = $"/subscriptions/{subscriptionId}"
                });
            
            var roleAssignment = new Assignment(
                "RoleAssignment", 
                new AssignmentArgs
                {
                    RoleDefinitionName = role.Name,
                    PrincipalId = azureADGroup.Id,
                    Scope = $"/subscriptions/{subscriptionId}"
                });

            const string resourceGroupName = "Pulumi";
            var resourceGroup = new ResourceGroup(
                "ResourceGroup",
                new ResourceGroupArgs
                {
                    Location = "CanadaEast",
                    Name = resourceGroupName
                });
            
            
            var application = new Application(
                "Application",
                new ApplicationArgs
                {
                    Name = "PulumiAzureSetup"
                });
            
            var servicePrincipal = new ServicePrincipal(
                "ServicePrincipal",
                new ServicePrincipalArgs
                {
                    ApplicationId = application.ApplicationId
                });
            
            var groupMembers = new GroupMember(
                "GroupMembership",
                new GroupMemberArgs
                {
                    GroupObjectId = azureADGroup.ObjectId,
                    MemberObjectId = servicePrincipal.ObjectId
                });
            
            var roleAssignmentResourceGroup = new Assignment(
                "RoleAssignmentRG", 
                new AssignmentArgs
                {
                    RoleDefinitionName = "Contributor",
                    PrincipalId = servicePrincipal.Id,
                    Scope = $"/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}"
                });
        });
    }
}
