package condition

import "ifttt/manager/common"

func (c *Condition) Manipulate(dependencies map[common.IntIota]any) error {
	if c.Group {
		for _, cnd := range c.Conditions {
			if err := cnd.Manipulate(dependencies); err != nil {
				return err
			}
		}
	} else if err := c.Operator1.Manipulate(dependencies); err != nil {
		return err
	} else if err := c.Operator2.Manipulate(dependencies); err != nil {
		return err
	}
	return nil
}
