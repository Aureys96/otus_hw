package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type Pipeline struct {
	stages []Stage
	input  In
	done   In
	output Out
}

func NewPipeline(input In, done In, stages []Stage) *Pipeline {
	return &Pipeline{stages: stages, input: input, done: done}
}

func (p *Pipeline) run() Out {
	for _, stage := range p.stages {
		intermediateChannel := make(Bi)

		go func(out Bi, in In) {
			defer close(out)

			for {
				select {
				case <-p.done:
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					out <- data
				}
			}
		}(intermediateChannel, p.input)

		p.input = stage(intermediateChannel)
	}

	p.output = p.input

	return p.output
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	pipeline := NewPipeline(in, done, stages)
	return pipeline.run()
}
