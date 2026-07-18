# Retention

Logs live for two hours.

- At the start of every run — before the wrapped command executes — swallow
  prunes the whole swallow dir tree: every log file whose modification time is
  older than two hours is deleted.
- Origin directories left empty by pruning are removed; the swallow dir
  itself never is.
- Pruning is best-effort: any failure is ignored and never affects the run.
- A log being written keeps a fresh modification time with every line, so
  active runs are safe from pruning. A run idle for over two hours may lose
  its log to a concurrent run's pruning — accepted.
- There is no daemon and no extra command; retention rides along with normal
  use.
