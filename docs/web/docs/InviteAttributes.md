# InviteAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | **string** | status of the invite | 
**Role** | **string** | role of the user in this city | 
**DistributorId** | [**uuid.UUID**](uuid.UUID.md) | distributor id | 
**Token** | **string** | unique token for the invite | 
**ExpiresAt** | **time.Time** | timestamp when the invite will expire | 
**CreatedAt** | **time.Time** | timestamp when the invite was created | 

## Methods

### NewInviteAttributes

`func NewInviteAttributes(status string, role string, distributorId uuid.UUID, token string, expiresAt time.Time, createdAt time.Time, ) *InviteAttributes`

NewInviteAttributes instantiates a new InviteAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInviteAttributesWithDefaults

`func NewInviteAttributesWithDefaults() *InviteAttributes`

NewInviteAttributesWithDefaults instantiates a new InviteAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *InviteAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InviteAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InviteAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetRole

`func (o *InviteAttributes) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *InviteAttributes) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *InviteAttributes) SetRole(v string)`

SetRole sets Role field to given value.


### GetDistributorId

`func (o *InviteAttributes) GetDistributorId() uuid.UUID`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *InviteAttributes) GetDistributorIdOk() (*uuid.UUID, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *InviteAttributes) SetDistributorId(v uuid.UUID)`

SetDistributorId sets DistributorId field to given value.


### GetToken

`func (o *InviteAttributes) GetToken() string`

GetToken returns the Token field if non-nil, zero value otherwise.

### GetTokenOk

`func (o *InviteAttributes) GetTokenOk() (*string, bool)`

GetTokenOk returns a tuple with the Token field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetToken

`func (o *InviteAttributes) SetToken(v string)`

SetToken sets Token field to given value.


### GetExpiresAt

`func (o *InviteAttributes) GetExpiresAt() time.Time`

GetExpiresAt returns the ExpiresAt field if non-nil, zero value otherwise.

### GetExpiresAtOk

`func (o *InviteAttributes) GetExpiresAtOk() (*time.Time, bool)`

GetExpiresAtOk returns a tuple with the ExpiresAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresAt

`func (o *InviteAttributes) SetExpiresAt(v time.Time)`

SetExpiresAt sets ExpiresAt field to given value.


### GetCreatedAt

`func (o *InviteAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *InviteAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *InviteAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


