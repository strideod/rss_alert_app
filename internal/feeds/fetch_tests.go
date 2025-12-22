// internal/feeds/fetch_test.go
package feeds

import (
    "bytes"
    "database/sql"
    "errors"
    "fmt"
    "html"
    "os"
    "testing"
    "time"

    "github.com/mmcdole/gofeed"
    "github.com/stretchr/testify/mock"
    "rss_alert_app/internal/db"
    "rss_alert_app/internal/models"
)

// Mocks for db functions
type MockDB struct {
    mock.Mock
}

func (m *MockDB) GetFeedURLs(sqlDB *sql.DB) ([]models.Feed, error) {
    args := m.Called(sqlDB)
    return args.Get(0).([]models.Feed), args.Error(1)
}
func (m *MockDB) IsSeen(sqlDB *sql.DB, url, guid string) (bool, error) {
    args := m.Called(sqlDB, url, guid)
    return args.Bool(0), args.Error(1)
}
func (m *MockDB) GetIncidentLastHash(sqlDB *sql.DB, feedID int, eventKey string) (string, bool, error) {
    args := m.Called(sqlDB, feedID, eventKey)
    return args.String(0), args.Bool(1), args.Error(2)
}
func (m *MockDB) InsertEventIgnore(sqlDB *sql.DB, ev models.Event) error {
    args := m.Called(sqlDB, ev)
    return args.Error(0)
}
func (m *MockDB) UpsertIncident(sqlDB *sql.DB, inc models.Incident) error {
    args := m.Called(sqlDB, inc)
    return args.Error(0)
}
func (m *MockDB) MarkSeen(sqlDB *sql.DB, url, guid string) error {
    args := m.Called(sqlDB, url, guid)
    return args.Error(0)
}

// Mocks for gofeed.Parser
type MockParser struct {
    mock.Mock
}

func (m *MockParser) ParseURL(url string) (*gofeed.Feed, error) {
    args := m.Called(url)
    return args.Get(0).(*gofeed.Feed), args.Error(1)
}

// Patchable functions
var (
    origGetFeedURLs      = db.GetFeedURLs
    origIsSeen           = db.IsSeen
    origGetIncidentLastHash = db.GetIncidentLastHash
    origInsertEventIgnore   = db.InsertEventIgnore
    origUpsertIncident      = db.UpsertIncident
    origMarkSeen            = db.MarkSeen
    origNewParser           = gofeed.NewParser
)

func patchDB(mockDB *MockDB) {
    db.GetFeedURLs = mockDB.GetFeedURLs
    db.IsSeen = mockDB.IsSeen
    db.GetIncidentLastHash = mockDB.GetIncidentLastHash
    db.InsertEventIgnore = mockDB.InsertEventIgnore
    db.UpsertIncident = mockDB.UpsertIncident
    db.MarkSeen = mockDB.MarkSeen
}
func restoreDB() {
    db.GetFeedURLs = origGetFeedURLs
    db.IsSeen = origIsSeen
    db.GetIncidentLastHash = origGetIncidentLastHash
    db.InsertEventIgnore = origInsertEventIgnore
    db.UpsertIncident = origUpsertIncident
    db.MarkSeen = origMarkSeen
}

func patchParser(mockParser *MockParser) {
    gofeed.NewParser = func() *gofeed.Parser {
        return &gofeed.Parser{
            ParseURL: mockParser.ParseURL,
        }
    }
}
func restoreParser() {
    gofeed.NewParser = origNewParser
}

func TestFetchFeeds_NoFeeds(t *testing.T) {
    mockDB := new(MockDB)
    patchDB(mockDB)
    defer restoreDB()

    mockDB.On("GetFeedURLs", mock.Anything).Return([]models.Feed{}, nil)

    // Capture output
    var buf bytes.Buffer
    stdout := os.Stdout
    os.Stdout = &buf
    defer func() { os.Stdout = stdout }()

    FetchFeeds(nil)

    output := buf.String()
    if !contains(output, "No feeds found") {
        t.Errorf("Expected warning about no feeds, got: %s", output)
    }
}

func TestFetchFeeds_ParseError(t *testing.T) {
    mockDB := new(MockDB)
    mockParser := new(MockParser)
    patchDB(mockDB)
    defer restoreDB()
    patchParser(mockParser)
    defer restoreParser()

    feeds := []models.Feed{{ID: 1, Name: "TestFeed", URL: "http://test.com/rss"}}
    mockDB.On("GetFeedURLs", mock.Anything).Return(feeds, nil)
    mockParser.On("ParseURL", "http://test.com/rss").Return((*gofeed.Feed)(nil), errors.New("parse error"))

    var buf bytes.Buffer
    stdout := os.Stdout
    os.Stdout = &buf
    defer func() { os.Stdout = stdout }()

    FetchFeeds(nil)

    output := buf.String()
    if !contains(output, "Failed to parse feed") {
        t.Errorf("Expected parse error log, got: %s", output)
    }
}

func TestFetchFeeds_NewEvent(t *testing.T) {
    mockDB := new(MockDB)
    mockParser := new(MockParser)
    patchDB(mockDB)
    defer restoreDB()
    patchParser(mockParser)
    defer restoreParser()

    feeds := []models.Feed{{ID: 1, Name: "TestFeed", URL: "http://test.com/rss"}}
    mockDB.On("GetFeedURLs", mock.Anything).Return(feeds, nil)

    item := &gofeed.Item{
        GUID:        "guid1",
        Title:       "Incident 1",
        Link:        "http://test.com/incidents/1",
        Description: "open",
        PublishedParsed: &time.Time{},
    }
    gf := &gofeed.Feed{Items: []*gofeed.Item{item}}
    mockParser.On("ParseURL", "http://test.com/rss").Return(gf, nil)

    mockDB.On("IsSeen", mock.Anything, "http://test.com/rss", "guid1").Return(false, nil)
    mockDB.On("GetIncidentLastHash", mock.Anything, 1, mock.Anything).Return("", false, nil)
    mockDB.On("InsertEventIgnore", mock.Anything, mock.Anything).Return(nil)
    mockDB.On("UpsertIncident", mock.Anything, mock.Anything).Return(nil)
    mockDB.On("MarkSeen", mock.Anything, "http://test.com/rss", "guid1").Return(nil)

    var buf bytes.Buffer
    stdout := os.Stdout
    os.Stdout = &buf
    defer func() { os.Stdout = stdout }()

    FetchFeeds(nil)

    output := buf.String()
    if !contains(output, "Incident 1") {
        t.Errorf("Expected event output, got: %s", output)
    }
}

func contains(s, substr string) bool {
    return bytes.Contains([]byte(s), []byte(substr))
}