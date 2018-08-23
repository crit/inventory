package dynamodb

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

var config = map[string]string{}

func init() {
	err := gotenv.Load(".testing")

	if err != nil {
		panic(err.Error())
	}

	config["region"] = os.Getenv("AWS_REGION")
	config["table"] = os.Getenv("GROWTV_DYNAMODB_TABLE")
}

func TestDynamoDB(t *testing.T) {
	store, err := New(config["region"], config["table"])

	assert.Nil(t, err)
	assert.NotNil(t, store)
}

type dynamoTester struct {
	key   string
	data  []byte
	count int64
	total int64
}

func (t *dynamoTester) Key() string {
	return t.key
}

func (t *dynamoTester) Type() string {
	return "storage.dynamoTester"
}

func (t *dynamoTester) Data() []byte {
	return t.data
}

func (t *dynamoTester) SetData(data []byte) {
	t.data = data
}

func TestDynamoDB_Read_404(t *testing.T) {
	str := dynamoTester{
		key: xid.New().String(),
	}

	store, err := New(config["region"], config["table"])
	assert.Nil(t, err)

	err = store.Read(&str)
	assert.EqualError(t, err, "[404] - not found")
}

func TestDynamoDB_WriteRead(t *testing.T) {
	payload, _ := json.Marshal(map[string]string{
		"alpha": "Alpha Value",
		"beta":  xid.New().String(),
	})

	a := dynamoTester{
		key:  xid.New().String(),
		data: payload,
	}

	store, err := New(config["region"], config["table"])
	assert.Nil(t, err)

	err = store.Write(&a)
	assert.Nil(t, err)

	b := dynamoTester{
		key: a.key,
	}

	err = store.Read(&b)

	assert.Nil(t, err)
	assert.Equal(t, payload, b.data)

	result := map[string]string{}
	err = json.Unmarshal(b.data, &result)

	assert.Nil(t, err)
	assert.Equal(t, "Alpha Value", result["alpha"])
}

func TestDynamoDB_WriteDelete(t *testing.T) {
	payload, _ := json.Marshal(map[string]string{
		"alpha": xid.New().String(),
	})

	a := dynamoTester{
		key:  xid.New().String(),
		data: payload,
	}

	store, err := New(config["region"], config["table"])
	assert.Nil(t, err)

	err = store.Write(&a)
	assert.Nil(t, err)

	err = store.Delete(&a)
	assert.Nil(t, err)
}

type dynamoTestLister struct {
	count int64
	total int64
	items []dynamoTester
}

func (d *dynamoTestLister) SetCount(count int64) {
	d.count = count
}

func (d *dynamoTestLister) Append(id, typ string, data []byte) {
	d.items = append(d.items, dynamoTester{key: id, data: data})
}

func (d *dynamoTestLister) Type() string {
	return "storage.dynamoTester"
}

func TestDynamoDB_List(t *testing.T) {
	a := dynamoTestLister{}

	store, err := New(config["region"], config["table"])
	assert.Nil(t, err)

	err = store.List(&a)
	assert.Nil(t, err)
	assert.Equal(t, a.count, int64(len(a.items)), "count != items")
}
