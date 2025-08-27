# DistributorAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Icon** | **string** | URL to the distributor&#39;s icon image. | 
**Name** | **string** | The name of the distributor. | 
**Status** | **string** | The current status of the distributor (e.g., active, inactive). | 
**UpdatedAt** | **time.Time** | The timestamp of the last update to the distributor&#39;s information. | 
**CreatedAt** | **time.Time** |  | 

## Methods

### NewDistributorAttributes

`func NewDistributorAttributes(icon string, name string, status string, updatedAt time.Time, createdAt time.Time, ) *DistributorAttributes`

NewDistributorAttributes instantiates a new DistributorAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDistributorAttributesWithDefaults

`func NewDistributorAttributesWithDefaults() *DistributorAttributes`

NewDistributorAttributesWithDefaults instantiates a new DistributorAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIcon

`func (o *DistributorAttributes) GetIcon() string`

GetIcon returns the Icon field if non-nil, zero value otherwise.

### GetIconOk

`func (o *DistributorAttributes) GetIconOk() (*string, bool)`

GetIconOk returns a tuple with the Icon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcon

`func (o *DistributorAttributes) SetIcon(v string)`

SetIcon sets Icon field to given value.


### GetName

`func (o *DistributorAttributes) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DistributorAttributes) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DistributorAttributes) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *DistributorAttributes) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DistributorAttributes) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DistributorAttributes) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetUpdatedAt

`func (o *DistributorAttributes) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *DistributorAttributes) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *DistributorAttributes) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetCreatedAt

`func (o *DistributorAttributes) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *DistributorAttributes) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *DistributorAttributes) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


