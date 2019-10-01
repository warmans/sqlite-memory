package sqlite

import (
	"testing"

	"github.com/go-joe/joe"
	"github.com/stretchr/testify/require"
)

func withInMemoryDb(t *testing.T, f func(mem joe.Memory)) {

	mem, err := NewMemory(":memory:")
	require.NoError(t, err)
	defer mem.Close()

	f(mem)
}

func TestMemory_Set(t *testing.T) {
	withInMemoryDb(t, func(mem joe.Memory) {
		//set a value
		err := mem.Set("foo", []byte("bar"))
		require.NoError(t, err)
	})
}

func TestMemory_Get(t *testing.T) {
	withInMemoryDb(t, func(mem joe.Memory) {
		// empty value
		val, found, err := mem.Get("foo")
		require.Nil(t, val)
		require.False(t, found)
		require.NoError(t, err)

		//set a value
		err = mem.Set("foo", []byte("bar"))
		require.NoError(t, err)

		// value found
		val, found, err = mem.Get("foo")
		require.EqualValues(t, []byte("bar"), val)
		require.True(t, found)
		require.NoError(t, err)
	})
}

func TestMemory_Delete(t *testing.T) {
	withInMemoryDb(t, func(mem joe.Memory) {

		//set a value
		err := mem.Set("foo", []byte("bar"))
		require.NoError(t, err)

		found, err := mem.Delete("foo")
		require.NoError(t, err)
		require.True(t, found)

		// value is gone
		val, found, err := mem.Get("foo")
		require.Nil(t, val)
		require.False(t, found)
		require.NoError(t, err)
	})
}

func TestMemory_Delete_NoneAffected(t *testing.T) {
	withInMemoryDb(t, func(mem joe.Memory) {

		ok, err := mem.Delete("foo")
		require.NoError(t, err)
		require.False(t, ok)
	})
}

func TestMemory_Keys(t *testing.T) {
	withInMemoryDb(t, func(mem joe.Memory) {
		keys := []string{"foo1", "foo2", "foo3"}

		for _, k := range keys {
			require.NoError(t, mem.Set(k, []byte(k+" value")))
		}

		foundKeys, err := mem.Keys()
		require.NoError(t, err)
		require.EqualValues(t, keys, foundKeys)

		for _, k := range foundKeys {
			v, ok, err := mem.Get(k)
			require.True(t, ok)
			require.NoError(t, err)
			require.EqualValues(t, []byte(k+" value"), v)
		}
	})
}
