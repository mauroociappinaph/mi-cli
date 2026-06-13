# API Reference

Generated on viernes, 12 de junio de 2026, 23:09:41 -03

## `cmd/sdd.go`
- func init() {

## `cmd/learning.go`
- func init() {

## `cmd/greet.go`
- func init() {

## `cmd/root.go`
- func Execute() error {
- func init() {
- func initConfig() error {

## `internal/engram/client.go`
- type Client struct {
- type Observation struct {
- type SearchResult struct {
- type SearchOptions struct {
- func NewClient() (*Client, error) {
- func (c *Client) initSchema() error {
- func (c *Client) Save(ctx context.Context, obs *Observation) (int64, error) {
- func (c *Client) Update(ctx context.Context, obs *Observation) error {
- func (c *Client) SaveOrUpdate(ctx context.Context, obs *Observation) (int64, error) {
- func (c *Client) Get(ctx context.Context, id int64) (*Observation, error) {

## `internal/learning/agent.go`
- type Agent struct {
- type Pattern struct {
- func NewAgent(scope string) (*Agent, error) {
- func (a *Agent) Close() error {
- func (a *Agent) Learn(ctx context.Context, pattern *Pattern) error {
- func (a *Agent) Recall(ctx context.Context, query string, limit int) ([]Pattern, error) {
- func (a *Agent) RecallByCategory(ctx context.Context, category string, limit int) ([]Pattern, error) {
- func (a *Agent) GetRecentPatterns(ctx context.Context, limit int) ([]Pattern, error) {
- func (a *Agent) parsePattern(content string) *Pattern {
- func (a *Agent) LearnFromSDD(ctx context.Context, phase, decision, rationale, files string) error {

