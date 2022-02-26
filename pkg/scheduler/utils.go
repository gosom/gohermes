package scheduler

import "sync"

func waitErrors(errs ...<-chan error) error {
	for err := range combineErrors(errs...) {
		if err != nil {
			return err
		}
	}
	return nil
}

func combineErrors(cs ...<-chan error) <-chan error {
	out := make(chan error, len(cs))
	var wg sync.WaitGroup

	output := func(c <-chan error) {
		defer wg.Done()
		for e := range c {
			out <- e
		}
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
