package todos

import "go.uber.org/dig"

// Module provides the todos module dependencies to the DI container.
func Module(c *dig.Container) error {
	if err := c.Provide(NewTodoRepository); err != nil {
		return err
	}

	if err := c.Provide(NewTodoService); err != nil {
		return err
	}

	if err := c.Provide(NewTodoHandler); err != nil {
		return err
	}

	return nil
}
