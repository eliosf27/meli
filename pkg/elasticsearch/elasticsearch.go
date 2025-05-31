package elasticsearch

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// Config holds the configuration for Elasticsearch.
type Config struct {
	Addresses []string
	Username  string
	Password  string
	// TODO: Add other relevant configuration fields like APIKey, CloudID, etc.
}

// Client is a wrapper around the Elasticsearch client.
type Client struct {
	es *elasticsearch.Client
	l  *log.Logger
}

// NewClient creates a new Elasticsearch client.
func NewClient(cfg Config, l *log.Logger) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		// TODO: Map other relevant fields from cfg to esCfg
		// For example, if you add APIKey or CloudID to your Config struct,
		// you would map them here.
		// APIKey: cfg.APIKey,
		// CloudID: cfg.CloudID,
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	return &Client{es: es, l: l}, nil
}

// Ping checks the connection to Elasticsearch.
func (c *Client) Ping() error {
	res, err := c.es.Ping()
	if err != nil {
		c.l.Printf("Error pinging Elasticsearch: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		c.l.Printf("Error pinging Elasticsearch: %s", res.Status())
		return fmt.Errorf("error pinging Elasticsearch: %s", res.Status())
	}

	c.l.Println("Successfully pinged Elasticsearch")
	return nil
}

// IndexDocument indexes a document in Elasticsearch.
// doc can be any type that can be marshalled to JSON (e.g., a struct or map[string]interface{}).
// docID is optional; if empty, Elasticsearch will generate an ID.
func (c *Client) IndexDocument(ctx context.Context, indexName string, docID string, doc interface{}) error {
	var body strings.Builder
	// TODO: This is a naive way to convert doc to JSON string.
	// For production, use a proper JSON marshaller.
	// For example, json.NewEncoder(&body).Encode(doc) and handle error.
	// This also assumes doc is already in a format that can be directly converted to string for simplicity.
	// A better approach would be to take []byte or io.Reader for the body.
	body.WriteString(fmt.Sprintf("%v", doc)) // Placeholder for proper JSON marshalling

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       strings.NewReader(body.String()),
		Refresh:    "true", // Or "wait_for" or "false" depending on consistency needs
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		c.l.Printf("Error indexing document %s in index %s: %s", docID, indexName, err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		c.l.Printf("Error indexing document %s in index %s: %s", docID, indexName, res.Status())
		return fmt.Errorf("error indexing document %s in index %s: %s", docID, indexName, res.Status())
	}

	c.l.Printf("Successfully indexed document %s in index %s", docID, indexName)
	return nil
}

// SearchDocuments searches for documents in Elasticsearch.
// It returns the search response or an error.
// For simplicity, this example takes a query string.
// A more robust implementation would take a structured query (e.g., map[string]interface{} or a custom query builder).
// The response body needs to be closed by the caller if no error.
func (c *Client) SearchDocuments(ctx context.Context, indexName string, query string) (*esapi.Response, error) {
	// For a more complex query, you would build a JSON query body.
	// Example: var queryBody strings.Builder
	// queryBody.WriteString(`{"query": {"match": {"title": "}})
	// queryBody.WriteString(query)
	// queryBody.WriteString(`"}}}`)
	// req.Body = strings.NewReader(queryBody.String())

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Query: query, // Uses the q= URI parameter for simplicity
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		c.l.Printf("Error searching in index %s with query '%s': %s", indexName, query, err)
		return nil, err
	}

	// Note: res.Body must be closed by the caller if err is nil.
	// defer res.Body.Close() // This would close it too early for the caller to read.

	if res.IsError() {
		// It's good practice to try and read the error response body if possible,
		// as it might contain more details.
		// var errBody bytes.Buffer // Example
		// io.Copy(&errBody, res.Body)
		// res.Body.Close()
		c.l.Printf("Error searching in index %s with query '%s': %s", indexName, query, res.Status())
		return res, fmt.Errorf("error searching in index %s: %s", indexName, res.Status())
	}

	c.l.Printf("Successfully performed search in index %s with query '%s'", indexName, query)
	return res, nil
}

// TODO: Add other relevant methods like DeleteDocument, UpdateDocument, etc.
