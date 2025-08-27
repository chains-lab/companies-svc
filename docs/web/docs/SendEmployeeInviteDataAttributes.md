# SendEmployeeInviteDataAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DistributorId** | **string** | The UUID of the distributor to which the user is being invited. | 
**UserId** | **string** | The UUID of the user being invited. | 
**Role** | **string** | The role assigned to the invited user within the distributor. | 

## Methods

### NewSendEmployeeInviteDataAttributes

`func NewSendEmployeeInviteDataAttributes(distributorId string, userId string, role string, ) *SendEmployeeInviteDataAttributes`

NewSendEmployeeInviteDataAttributes instantiates a new SendEmployeeInviteDataAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSendEmployeeInviteDataAttributesWithDefaults

`func NewSendEmployeeInviteDataAttributesWithDefaults() *SendEmployeeInviteDataAttributes`

NewSendEmployeeInviteDataAttributesWithDefaults instantiates a new SendEmployeeInviteDataAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDistributorId

`func (o *SendEmployeeInviteDataAttributes) GetDistributorId() string`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *SendEmployeeInviteDataAttributes) GetDistributorIdOk() (*string, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *SendEmployeeInviteDataAttributes) SetDistributorId(v string)`

SetDistributorId sets DistributorId field to given value.


### GetUserId

`func (o *SendEmployeeInviteDataAttributes) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *SendEmployeeInviteDataAttributes) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *SendEmployeeInviteDataAttributes) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetRole

`func (o *SendEmployeeInviteDataAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *SendEmployeeInviteDataAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *SendEmployeeInviteDataAttributes) SetRole(v string)`

SetRole sets Role field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


