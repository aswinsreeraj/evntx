package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aswinsreeraj/evntx/internal/domain"
	"github.com/aswinsreeraj/evntx/pkg/jwt"
	"github.com/aswinsreeraj/evntx/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func setTestEnv() {
	os.Setenv("JWT_SECRET", "super-secret-key")
	os.Setenv("JWT_REFRESH_SECRET", "super-secret-refresh-key")
}

func getAuthHeader(userID string) string {
	token, _ := jwt.GenerateAccessToken(userID)
	return "Bearer " + token
}

func TestWorkflowE2E_HealthCheck(t *testing.T) {
	setTestEnv()
	db := testutil.SetupTestDB(t)
	defer testutil.ClearDatabase(db)

	router := testutil.SetupTestRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestWorkflowE2E_BookingJourney(t *testing.T) {
	setTestEnv()
	db := testutil.SetupTestDB(t)
	defer testutil.ClearDatabase(db)

	router := testutil.SetupTestRouter(db)

	
	user := testutil.SeedUser(db, "e2e_user@example.com", "user")
	organizer := testutil.SeedUser(db, "e2e_org@example.com", string(domain.RoleOrganizer))
	
	
	event := testutil.SeedEvent(db, organizer.ID, "E2E Music Festival", "live")
	ticketType := testutil.SeedTicketType(db, event.ID, "VIP", 500.0, 100)

	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/events/"+event.Slug, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "E2E Music Festival")

	
	reservePayload := map[string]interface{}{
		"event_id": event.ID,
		"tickets": []map[string]interface{}{
			{
				"ticket_type_id": ticketType.ID,
				"quantity":       2,
			},
		},
	}
	body, _ := json.Marshal(reservePayload)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/bookings/reserve", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", getAuthHeader(user.ID))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code, "Expected successful ticket reservation")

	var reserveResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &reserveResp)
	assert.NoError(t, err)
	
	var bookingID string
	if data, ok := reserveResp["data"].(map[string]interface{}); ok {
		bookingID = data["booking_id"].(string)
	} else if reserveResp["booking_id"] != nil {
		bookingID = reserveResp["booking_id"].(string)
	}
	assert.NotEmpty(t, bookingID)

	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/me/bookings?status=reserved", nil)
	req.Header.Set("Authorization", getAuthHeader(user.ID))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), bookingID)
}

func TestWorkflowE2E_AdminEventApproval(t *testing.T) {
	setTestEnv()
	db := testutil.SetupTestDB(t)
	defer testutil.ClearDatabase(db)

	router := testutil.SetupTestRouter(db)

	
	admin := testutil.SeedUser(db, "e2e_admin@example.com", string(domain.RoleAdmin))
	organizer := testutil.SeedUser(db, "e2e_org2@example.com", string(domain.RoleOrganizer))
	
	
	event := testutil.SeedEvent(db, organizer.ID, "Pending Tech Talk", "pending")

	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/events?status=pending", nil)
	req.Header.Set("Authorization", getAuthHeader(admin.ID))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Pending Tech Talk")

	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PATCH", "/admin/events/"+event.ID+"/approve", nil)
	req.Header.Set("Authorization", getAuthHeader(admin.ID))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/events/"+event.Slug, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "Expected event to be live after admin approval")
}
