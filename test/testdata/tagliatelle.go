// args: -Etagliatelle
package testdata

import "time"

type TglFoo struct {
	ID     string        `json:"ID"`     // ERROR `json\(camel\): got 'ID' want 'id'`
	UserID string        `json:"UserID"` // ERROR `json\(camel\): got 'UserID' want 'userId'`
	Name   string        `json:"name"`
	Value  time.Duration `json:"value,omitempty"`
	Bar    TglBar        `json:"bar"`
	Bur    `json:"bur"`
}

type TglBar struct {
	Name                 string  `json:"-"`
	Value                string  `json:"value"`
	CommonServiceFooItem *TglBir `json:"CommonServiceItem,omitempty"` // ERROR `json\(camel\): got 'CommonServiceItem' want 'commonServiceItem'`
}

type TglBir struct {
	Name             string   `json:"-"`
	Value            string   `json:"value"`
	ReplaceAllowList []string `mapstructure:"replace-allow-list"`
}

type Bur struct {
	Name    string
	Value   string `yaml:"Value"` // ERROR `yaml\(camel\): got 'Value' want 'value'`
	More    string `json:"-"`
	Also    string `json:"also,omitempty"`
	ReqPerS string `avro:"req_per_s"`
}
