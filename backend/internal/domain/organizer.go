package domain

type StatCard struct {
	Value		float64	`json:"value"`
	Percentage	float64	`json:"percentage"`
	Label		string	`json:"label,omitempty"`
}

type RevenuePoint struct {
	Date	string	`json:"date"`
	Amount	float64	`json:"amount"`
}

type EventSalesBreakdown struct {
	EventName	string	`json:"name"`
	Revenue		float64	`json:"value"`
}

type OrganizerDashboardStats struct {
	TotalRevenue	StatCard		`json:"total_revenue"`
	TicketsSold	StatCard		`json:"tickets_sold"`
	ActiveEvents	StatCard		`json:"active_events"`
	PendingEvents	StatCard		`json:"pending_events"`
	RevenueOverview	[]RevenuePoint		`json:"revenue_overview"`
	SalesBreakdown	[]EventSalesBreakdown	`json:"sales_breakdown"`
}

type TicketSalesProportion struct {
	EventName	string	`json:"name"`
	TicketsSold	int	`json:"tickets_sold"`
	PercentageTotal	float64	`json:"percentage_total"`
}

type SalesReportStats struct {
	TotalRevenue	StatCard		`json:"total_revenue"`
	TicketsSold	StatCard		`json:"tickets_sold"`
	RevenueOverTime	[]RevenuePoint		`json:"revenue_over_time"`
	TicketsPerEvent	[]TicketSalesProportion	`json:"tickets_per_event"`
}

type FunnelStep struct {
	Label		string	`json:"label"`
	Count		int	`json:"count"`
	Percentage	float64	`json:"percentage"`
}

type PeakUsagePoint struct {
	Label		string	`json:"label"`
	Viewing		int	`json:"viewing"`
	Bookings	int	`json:"bookings"`
}

type EngagementReportStats struct {
	PageViews	StatCard		`json:"page_views"`
	ConversionRate	StatCard		`json:"conversion_rate"`
	UserJourney	[]FunnelStep		`json:"user_journey"`
	PeakUsage	[]PeakUsagePoint	`json:"peak_usage"`
}

type AdminStatCard struct {
	Value		float64	`json:"value"`
	Percentage	float64	`json:"percentage"`
	Subtitle	string	`json:"subtitle,omitempty"`
}

type AdminDashboardStats struct {
	Revenue			AdminStatCard	`json:"revenue"`
	TotalUsers		AdminStatCard	`json:"total_users"`
	TotalOrganizers		AdminStatCard	`json:"total_organizers"`
	TotalEvents		AdminStatCard	`json:"total_events"`
	RefundRate		AdminStatCard	`json:"refund_rate"`
	UserGrowth		AdminStatCard	`json:"user_growth"`
	PendingApprovals	AdminStatCard	`json:"pending_approvals"`
	ActiveEvents		AdminStatCard	`json:"active_events"`
	RevenueOverview		[]RevenuePoint	`json:"revenue_overview"`
}

type CategoryRevenue struct {
	Category	string	`json:"category"`
	Revenue		float64	`json:"revenue"`
}

type RefundDataPoint struct {
	Month	string	`json:"month"`
	Amount	float64	`json:"amount"`
}

type TopOrganizerEntry struct {
	Name		string	`json:"name"`
	Revenue		float64	`json:"revenue"`
	ActiveEvents	int	`json:"active_events"`
	PendingEvents	int	`json:"pending_events"`
	AvgEventRating	float64	`json:"avg_event_rating"`
}

type TopUserEntry struct {
	Name		string	`json:"name"`
	EventsAttended	int	`json:"events_attended"`
	TotalSpent	float64	`json:"total_spent"`
}

type AdminRevenueReport struct {
	RevenueToday		AdminStatCard		`json:"revenue_today"`
	RevenueThisMonth	AdminStatCard		`json:"revenue_this_month"`
	TotalRevenue		AdminStatCard		`json:"total_revenue"`
	GrowthRate		AdminStatCard		`json:"growth_rate"`
	RevenueOverTime		[]RevenuePoint		`json:"revenue_over_time"`
	CategoryBreakdown	[]CategoryRevenue	`json:"category_breakdown"`
	RefundAnalytics		[]RefundDataPoint	`json:"refund_analytics"`
	RefundTotal		AdminStatCard		`json:"refund_total"`
	TopOrganizers		[]TopOrganizerEntry	`json:"top_organizers"`
	TopUsers		[]TopUserEntry		`json:"top_users"`
}
