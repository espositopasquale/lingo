package rules

type Engine struct {
	rules []Rule
}

func NewEngine(rules ...Rule) *Engine {
	return &Engine{
		rules: append([]Rule(nil), rules...),
	}
}

func (e *Engine) Match(text string, start int) (Match, bool) {
	var best Match
	bestPriority := 0
	found := false

	for _, rule := range e.rules {
		match, ok := rule.Match(text, start)

		if !ok {
			continue
		}

		if match.Start != start {
			continue
		}

		if match.End <= match.Start {
			continue
		}

		if match.End > len(text) {
			continue
		}

		priority := rule.Priority()

		if !found ||
			priority > bestPriority ||
			(priority == bestPriority &&
				match.End > best.End) {

			best = match
			bestPriority = priority
			found = true
		}
	}

	return best, found
}
