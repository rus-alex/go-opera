package opera

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func UpdateRules(src Rules, diff []byte) (res Rules, err error) {
	changed := src
	err = json.Unmarshal(diff, &changed)
	if err != nil {
		return
	}
	// protect readonly fields
	res = src
	res.Blocks = changed.Blocks
	res.Dag = changed.Dag
	res.Economy = changed.Economy
	return
}

func TestUpdateRules(t *testing.T) {
	require := require.New(t)

	var exp Rules
	exp.Dag.MaxEpochBlocks = 99

	exp.Dag.MaxParents = 5
	exp.Economy.BlockMissedSlack = 7
	exp.Blocks.BlockGasHardLimit = 1000
	got, err := UpdateRules(exp, []byte(`{"Dag":{"MaxParents":5},"Economy":{"BlockMissedSlack":7},"Blocks":{"BlockGasHardLimit":1000}}`))
	require.NoError(err)
	require.Equal(exp, got, "mutate fields")

	exp.Dag.MaxParents = 0
	got, err = UpdateRules(exp, []byte(`{"Name":"xxx","NetworkID":1,"Dag":{"MaxParents":0}}`))
	require.NoError(err)
	require.Equal(exp, got, "readonly fields")

	got, err = UpdateRules(exp, []byte(`{}`))
	require.NoError(err)
	require.Equal(exp, got, "empty diff")

	_, err = UpdateRules(exp, []byte(`}{`))
	require.Error(err)
}
