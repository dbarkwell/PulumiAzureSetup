import pulumi
import pulumi_azuread
from pulumi_azure import ad, core, role

config = pulumi.Config()
subscription_id = config.require("subscriptionId")

azure_ad_group = pulumi_azuread.Group("AzureUsersGroup", name="All Azure Users")

scope = f"/subscriptions/{subscription_id}"
role_definition = role.Definition("Role",
                                  name="Register Azure resource providers",
                                  description="Can register Azure resource providers",
                                  permissions=[{"actions": ["*/register/action"]}],
                                  assignable_scopes=[scope],
                                  scope=scope)

role_assignment = role.Assignment("RoleAssignment",
                                  role_definition_name=role_definition.name,
                                  principal_id=azure_ad_group.id,
                                  scope=scope)

resource_group_name = "Pulumi"
resource_group = core.ResourceGroup("ResourceGroup",
                                    location="CanadaEast",
                                    name=resource_group_name)

application = ad.Application("Application",
                             name="PulumiAzureSetup")

service_principal = ad.ServicePrincipal("ServicePrincipal",
                                        application_id=application.application_id)

group_members = pulumi_azuread.GroupMember("GroupMembership",
                                           group_object_id=azure_ad_group.object_id,
                                           member_object_id=service_principal.id)

role_assignment_resource_group = role.Assignment("RoleAssignmentRG",
                                                 role_definition_name="Contributor",
                                                 principal_id=service_principal.id,
                                                 scope=f"{scope}/resourceGroups/{resource_group_name}")
