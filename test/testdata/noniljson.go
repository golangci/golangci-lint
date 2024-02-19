//golangcitest:args -Enoniljson
package testdata

type userData struct {
	Name     *string           `json:"name"` // want `nullable field 'Name' in struct 'userData' must include 'omitempty' in its json tag to avoid marshaling as null`
	Email    string            `json:"email,omitempty"`
	Address  *string           `json:"address,omitempty"`
	Friends  []string          `json:"friends"`  // want `nullable field 'Friends' in struct 'userData' must include 'omitempty' in its json tag to avoid marshaling as null`
	Metadata map[string]string `json:"metadata"` // want `nullable field 'Metadata' in struct 'userData' must include 'omitempty' in its json tag to avoid marshaling as null`
}

type dynamicData struct {
	Data interface{} `json:"data"` // want `nullable field 'Data' in struct 'dynamicData' must include 'omitempty' in its json tag to avoid marshaling as null`
}

type productData struct {
	ID           int                    `json:"id"`
	Description  *string                `json:"description"` // want `nullable field 'Description' in struct 'productData' must include 'omitempty' in its json tag to avoid marshaling as null`
	Price        float64                `json:"price,omitempty"`
	Availability []int                  `json:"availability"` // want `nullable field 'Availability' in struct 'productData' must include 'omitempty' in its json tag to avoid marshaling as null`
	Attributes   map[string]interface{} `json:"attributes"`   // want `nullable field 'Attributes' in struct 'productData' must include 'omitempty' in its json tag to avoid marshaling as null`
}
