package rule

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"
)

func (rs *RuleSwitch) Manipulate(dependencies map[common.IntIota]any) error {
	for idx, c := range rs.Cases {
		if err := c.Manipulate(false, dependencies); err != nil {
			return err
		} else {
			rs.Cases[idx] = c
		}
	}
	if err := rs.Default.Manipulate(true, dependencies); err != nil {
		return err
	}
	return nil
}

func (rsc *RuleSwitchCase) Manipulate(isDefault bool, dependencies map[common.IntIota]any) error {
	if !isDefault {
		if err := rsc.Condition.Manipulate(dependencies); err != nil {
			return err
		}
	}
	if manipulated, err := resolvable.ManipulateArray(&rsc.Do, dependencies); err != nil {
		return err
	} else {
		rsc.Do = *manipulated
	}
	return nil
}
