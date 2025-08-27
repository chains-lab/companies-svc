# DistributorBlockAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DistributorId** | **string** | ID of the distributor being blocked | 
**InitiatorId** | **string** | ID of the user who initiated the block | 
**Reason** | **string** | Reason for blocking the distributor | 
**Status** | **string** | Current status of the block | 
**BlockedAt** | **time.Time** | Timestamp when the block was initiated | 
**CancelledAt** | Pointer to **time.Time** | Timestamp when the block was lifted, if applicable | [optional] 

## Methods

### NewDistributorBlockAttributes

`func NewDistributorBlockAttributes(distributorId string, initiatorId string, reason string, status string, blockedAt time.Time, ) *DistributorBlockAttributes`

NewDistributorBlockAttributes instantiates a new DistributorBlockAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDistributorBlockAttributesWithDefaults

`func NewDistributorBlockAttributesWithDefaults() *DistributorBlockAttributes`

NewDistributorBlockAttributesWithDefaults instantiates a new DistributorBlockAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDistributorId

`func (o *DistributorBlockAttributes) GetDistributorId() string`

GetDistributorId returns the DistributorId field if non-nil, zero value otherwise.

### GetDistributorIdOk

`func (o *DistributorBlockAttributes) GetDistributorIdOk() (*string, bool)`

GetDistributorIdOk returns a tuple with the DistributorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistributorId

`func (o *DistributorBlockAttributes) SetDistributorId(v string)`

SetDistributorId sets DistributorId field to given value.


### GetInitiatorId

`func (o *DistributorBlockAttributes) GetInitiatorId() string`

GetInitiatorId returns the InitiatorId field if non-nil, zero value otherwise.

### GetInitiatorIdOk

`func (o *DistributorBlockAttributes) GetInitiatorIdOk() (*string, bool)`

GetInitiatorIdOk returns a tuple with the InitiatorId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitiatorId

`func (o *DistributorBlockAttributes) SetInitiatorId(v string)`

SetInitiatorId sets InitiatorId field to given value.


### GetReason

`func (o *DistributorBlockAttributes) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *DistributorBlockAttributes) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *DistributorBlockAttributes) SetReason(v string)`

SetReason sets Reason field to given value.


### GetStatus

`func (o *DistributorBlockAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DistributorBlockAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DistributorBlockAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetBlockedAt

`func (o *DistributorBlockAttributes) GetBlockedAt() time.Time`

GetBlockedAt returns the BlockedAt field if non-nil, zero value otherwise.

### GetBlockedAtOk

`func (o *DistributorBlockAttributes) GetBlockedAtOk() (*time.Time, bool)`

GetBlockedAtOk returns a tuple with the BlockedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockedAt

`func (o *DistributorBlockAttributes) SetBlockedAt(v time.Time)`

SetBlockedAt sets BlockedAt field to given value.


### GetCancelledAt

`func (o *DistributorBlockAttributes) GetCancelledAt() time.Time`

GetCancelledAt returns the CancelledAt field if non-nil, zero value otherwise.

### GetCancelledAtOk

`func (o *DistributorBlockAttributes) GetCancelledAtOk() (*time.Time, bool)`

GetCancelledAtOk returns a tuple with the CancelledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCancelledAt

`func (o *DistributorBlockAttributes) SetCancelledAt(v time.Time)`

SetCancelledAt sets CancelledAt field to given value.

### HasCancelledAt

`func (o *DistributorBlockAttributes) HasCancelledAt() bool`

HasCancelledAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


