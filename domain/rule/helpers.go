package rule

import (
	"ifttt/manager/common"
	"ifttt/manager/domain/resolvable"
)

func (rs *RuleSwitch) Manipulate(dependencies map[common.IntIota]any) error {
	for _, c := range rs.Cases {
		if err := c.Manipulate(dependencies); err != nil {
			return err
		}
	}
	if err := rs.Default.Manipulate(dependencies); err != nil {
		return err
	}
	return nil
}

func (rsc *RuleSwitchCase) Manipulate(dependencies map[common.IntIota]any) error {
	if err := rsc.Condition.Manipulate(dependencies); err != nil {
		return err
	} else if manipulated, err := resolvable.ManipulateArray(rsc.Do, dependencies); err != nil {
		return err
	} else {
		rsc.Do = manipulated
	}
	return nil
}
