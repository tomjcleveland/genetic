package genetic

import "sort"

type result struct {
	result indWithScore
	err    error
}

func calculateFitnessConcurrently(in []indWithScore, workers int) (out []indWithScore, err error) {
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan indWithScore)
	results := make(chan result)

	// Spin up workers
	for i := 0; i < workers; i++ {
		go func() {
			for {
				job, ok := <-jobs
				if !ok {
					break
				}
				score, err := job.Individual.Fitness()
				job.score = score
				results <- result{result: job, err: err}
			}
		}()
	}

	// Add jobs
	go func() {
		for _, ind := range in {
			jobs <- ind
		}
		close(jobs)
	}()

	// Get results
	for i := 0; i < len(in); i++ {
		outcome := <-results
		if outcome.err != nil {
			return nil, err
		}
		out = append(out, outcome.result)
	}
	close(results)
	sort.Sort(sort.Reverse(pairs(out)))

	return out, nil
}
