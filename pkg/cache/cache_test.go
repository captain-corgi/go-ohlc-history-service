package cache

import (
	"reflect"
	"testing"
)

func TestInMemoryDBImpl_Clear(t *testing.T) {
	type fields struct {
		cache map[string]interface{}
	}
	tests := []struct {
		fields fields
		name   string
	}{
		{
			name: "Test InMemoryDBImpl Clear",
			fields: fields{
				cache: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &InMemoryDBImpl{
				cache: tt.fields.cache,
			}
			db.Clear()
		})
	}
}

func TestInMemoryDBImpl_Delete(t *testing.T) {
	type fields struct {
		cache map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test InMemoryDBImpl Delete",
			fields: fields{
				cache: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				key: "key1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &InMemoryDBImpl{
				cache: tt.fields.cache,
			}
			db.Delete(tt.args.key)
		})
	}
}

func TestInMemoryDBImpl_Get(t *testing.T) {
	type fields struct {
		cache map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		want   interface{}
		fields fields
		name   string
		args   args
		want1  bool
	}{
		{
			name: "Test InMemoryDBImpl Get",
			fields: fields{
				cache: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				key: "key1",
			},
			want:  "value1",
			want1: true,
		},
		{
			name: "Test InMemoryDBImpl Get",
			fields: fields{
				cache: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				key: "key3",
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &InMemoryDBImpl{
				cache: tt.fields.cache,
			}
			got, got1 := db.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInMemoryDBImpl_Set(t *testing.T) {
	type fields struct {
		cache map[string]interface{}
	}
	type args struct {
		value interface{}
		key   string
	}
	tests := []struct {
		args   args
		fields fields
		name   string
	}{
		{
			name: "Test InMemoryDBImpl Set",
			fields: fields{
				cache: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				key:   "key3",
				value: "value3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &InMemoryDBImpl{
				cache: tt.fields.cache,
			}
			db.Set(tt.args.key, tt.args.value)
		})
	}
}

func TestNewInMemoryDB(t *testing.T) {
	tests := []struct {
		want *InMemoryDBImpl
		name string
	}{
		{
			name: "NewInMemoryDB",
			want: &InMemoryDBImpl{
				cache: make(map[string]interface{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemoryDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemoryDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
