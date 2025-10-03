# CompanyBlockAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompanyId** | [**uuid.UUID**](uuid.UUID.md) | ID of the company being blocked | 
**InitiatorId** | [**uuid.UUID**](uuid.UUID.md) | ID of the user who initiated the block | 
**Reason** | **string** | Reason for blocking the company | 
**Status** | **string** | Current status of the block | 
**BlockedAt** | **time.Time** | Timestamp when the block was initiated | 
**CancelledAt** | Pointer to **time.Time** | Timestamp when the block was lifted, if applicable | [optional] 

## Methods

### NewCompanyBlockAttributes

`func NewCompanyBlockAttributes(companyId uuid.UUID, initiatorId uuid.UUID, reason string, status string, blockedAt time.Time, ) *CompanyBlockAttributes`

NewCompanyBlockAttributes instantiates a new CompanyBlockAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCompanyBlockAttributesWithDefaults

`func NewCompanyBlockAttributesWithDefaults() *CompanyBlockAttributes`

NewCompanyBlockAttributesWithDefaults instantiates a new CompanyBlockAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompanyId

`func (o *CompanyBlockAttributes) GetCompanyId() uuid.UUID`

GetCompanyId returns the CompanyId field if non-nil, zero value otherwise.

### GetCompanyIdOk

`func (o *CompanyBlockAttributes) GetCompanyIdOk() (*uuid.UUID, bool)`

GetCompanyIdOk returns a tuple with the CompanyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompanyId

`func (o *CompanyBlockAttributes) SetCompanyId(v uuid.UUID)`

SetCompanyId sets CompanyId field to given value.


### GetInitiatorId

`func (o *CompanyBlockAttributes) GetInitiatorId() uuid.UUID`

GetInitiatorId returns the InitiatorId field if non-nil, zero value otherwise.

### GetInitiatorIdOk

`func (o *CompanyBlockAttributes) GetInitiatorIdOk() (*uuid.UUID, bool)`

GetInitiatorIdOk returns a tuple with the InitiatorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitiatorId

`func (o *CompanyBlockAttributes) SetInitiatorId(v uuid.UUID)`

SetInitiatorId sets InitiatorId field to given value.


### GetReason

`func (o *CompanyBlockAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *CompanyBlockAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *CompanyBlockAttributes) SetReason(v string)`

SetReason sets Reason field to given value.


### GetStatus

`func (o *CompanyBlockAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CompanyBlockAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CompanyBlockAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetBlockedAt

`func (o *CompanyBlockAttributes) GetBlockedAt() time.Time`

GetBlockedAt returns the BlockedAt field if non-nil, zero value otherwise.

### GetBlockedAtOk

`func (o *CompanyBlockAttributes) GetBlockedAtOk() (*time.Time, bool)`

GetBlockedAtOk returns a tuple with the BlockedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockedAt

`func (o *CompanyBlockAttributes) SetBlockedAt(v time.Time)`

SetBlockedAt sets BlockedAt field to given value.


### GetCancelledAt

`func (o *CompanyBlockAttributes) GetCancelledAt() time.Time`

GetCancelledAt returns the CancelledAt field if non-nil, zero value otherwise.

### GetCancelledAtOk

`func (o *CompanyBlockAttributes) GetCancelledAtOk() (*time.Time, bool)`

GetCancelledAtOk returns a tuple with the CancelledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCancelledAt

`func (o *CompanyBlockAttributes) SetCancelledAt(v time.Time)`

SetCancelledAt sets CancelledAt field to given value.

### HasCancelledAt

`func (o *CompanyBlockAttributes) HasCancelledAt() bool`

HasCancelledAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


