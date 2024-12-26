package triggerflow

import "ifttt/manager/common"

func (tc *TriggerConditionRequest) Manipulate(dependencies map[common.IntIota]any) error {
	if err := tc.If.Manipulate(dependencies); err != nil {
		return err
	}
	return nil
}
