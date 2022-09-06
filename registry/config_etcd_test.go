package registry

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/client/v2"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Run("successfully get key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Get", "test/testKey").Return("testValue", nil)

		sut := &etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		value, err := Get(sut.parent, "testKey", sut.client)

		// then
		require.NoError(t, err)
		assert.Equal(t, "testValue", value)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})

	t.Run("error on getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Get", "test/testKey").Return("", assert.AnError)

		sut := &etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		_, err := Get(sut.parent, "testKey", sut.client)

		// then
		require.Error(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_Delete(t *testing.T) {
	t.Run("Successfully delete key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Delete", "test/testKey", mock.Anything).Return(nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.Delete("testKey")

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})

	t.Run("error on deleting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Delete", "test/testKey", mock.Anything).Return(assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.Delete("testKey")

		// then
		require.Error(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_DeleteRecursive(t *testing.T) {
	t.Run("Successfully delete key recursively from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("DeleteRecursive", "test/testKey").Return(nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.DeleteRecursive("testKey")

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})

	t.Run("error on deleting key recursively from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("DeleteRecursive", "test/testKey").Return(assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.DeleteRecursive("testKey")

		// then
		require.Error(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_Exists(t *testing.T) {
	t.Run("Successfully checking key existence in registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Exists", "test/testKey", mock.Anything).Return(true, nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		ok, err := sut.Exists("testKey")

		// then
		require.NoError(t, err)
		assert.True(t, ok)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})

	t.Run("error on checking key existence in registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Exists", "test/testKey", mock.Anything).Return(false, assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		ok, err := sut.Exists("testKey")

		// then
		require.Error(t, err, assert.AnError)
		assert.False(t, ok)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_Get(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Get", "test/testKey").Return("testValue", nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		value, err := sut.Get("testKey")

		// then
		require.NoError(t, err)
		assert.Equal(t, "testValue", value)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_GetAll(t *testing.T) {
	t.Run("Successfully get all keys in registry", func(t *testing.T) {
		// given
		testKeys := map[string]string{
			"test1Key": "test1Value",
			"test2Key": "test2Value",
		}

		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("GetRecursive", "test").Return(testKeys, nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		actualKeyMap, err := sut.GetAll()

		// then
		require.NoError(t, err)
		assert.Equal(t, testKeys["test1Key"], actualKeyMap["test1Key"])
		assert.Equal(t, testKeys["test2Key"], actualKeyMap["test2Key"])
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})

	t.Run("error on getting all keys in registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("GetRecursive", "test").Return(map[string]string{}, assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		_, err := sut.GetAll()

		// then
		require.Error(t, err, assert.AnError)
	})
}

func Test_etcdConfigurationContext_GetOrFalse(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Exists", "test/testKey").Return(true, nil)
		etcdClientMock.On("Get", "test/testKey").Return("testValue", nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		ok, value, err := sut.GetOrFalse("testKey")

		// then
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "testValue", value)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})

	t.Run("error on checking existence of key", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Exists", "test/testKey").Return(false, assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		ok, _, err := sut.GetOrFalse("testKey")

		// then
		require.Error(t, err, assert.AnError)
		assert.False(t, ok)
	})

	t.Run("return false if key does not exist", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Exists", "test/testKey").Return(false, nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		ok, _, err := sut.GetOrFalse("testKey")

		// then
		require.NoError(t, err, assert.AnError)
		assert.False(t, ok)
	})

	t.Run("error when getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Exists", "test/testKey").Return(true, nil)
		etcdClientMock.On("Get", "test/testKey").Return("", assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		ok, _, err := sut.GetOrFalse("testKey")

		// then
		require.Error(t, err, assert.AnError)
		assert.False(t, ok)
	})
}

func Test_etcdConfigurationContext_Refresh(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		options := &client.SetOptions{
			TTL:     time.Second * 500,
			Refresh: true,
		}

		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Set", "test/testKey", "", options).Return("", nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.Refresh("testKey", 500)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_RemoveAll(t *testing.T) {
	t.Run("successfully delete all keys recursively key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("DeleteRecursive", "test").Return(nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.RemoveAll()

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
	t.Run("error on deleting all keys recursively key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("DeleteRecursive", "test").Return(assert.AnError)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.RemoveAll()

		// then
		require.Error(t, err, assert.AnError)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_Set(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Set", "test/testKey", "testValue", mock.Anything).Return("", nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.Set("testKey", "testValue")

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdConfigurationContext_SetWithLifetime(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Set", "test/testKey", "testValue", &client.SetOptions{
			TTL: time.Second * 500,
		}).Return("", nil)

		sut := etcdConfigurationContext{
			parent: "test",
			client: etcdClientMock,
		}

		// when
		err := sut.SetWithLifetime("testKey", "testValue", 500)

		// then
		require.NoError(t, err)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdWatchConfigurationContext_Get(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Get", "/testKey").Return("testValue", nil)

		sut := etcdWatchConfigurationContext{
			client: etcdClientMock,
		}

		// when
		value, err := sut.Get("testKey")

		// then
		require.NoError(t, err)
		assert.Equal(t, "testValue", value)
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdWatchConfigurationContext_GetChildrenPaths(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("GetChildrenPaths", "testKey").Return([]string{"test1", "test2"}, nil)

		sut := etcdWatchConfigurationContext{
			client: etcdClientMock,
		}

		// when
		childrenPaths, err := sut.GetChildrenPaths("testKey")

		// then
		require.NoError(t, err)
		assert.Equal(t, "test1", childrenPaths[0])
		assert.Equal(t, "test2", childrenPaths[1])
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}

func Test_etcdWatchConfigurationContext_Watch(t *testing.T) {
	t.Run("Successfully getting key from registry", func(t *testing.T) {
		// given
		eventChannel := make(chan *client.Response)

		etcdClientMock := &etcdClientMock{}
		etcdClientMock.On("Watch", context.Background(), "testKey", true, eventChannel).Return()

		sut := etcdWatchConfigurationContext{
			client: etcdClientMock,
		}

		// when
		sut.Watch(context.Background(), "testKey", true, eventChannel)

		// then
		mock.AssertExpectationsForObjects(t, etcdClientMock)
	})
}
