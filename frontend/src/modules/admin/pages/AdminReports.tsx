import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import AdminLayout from "../components/AdminLayout";
import { adminApi } from "../api";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
} from "recharts";
import { ChevronDown, Loader2, Download } from "lucide-react";
import { exportToCSV } from "../../../shared/utils/csv";

const DATE_RANGES = [
  { label: "Last 30 Days", value: "30D" },
  { label: "Last 90 Days", value: "90D" },
  { label: "Last Year", value: "1Y" },
  { label: "All Time", value: "ALL" },
];

const CATEGORY_COLORS = [
  "#4F46E5", "#2DD4BF", "#FBBF24", "#60A5FA", "#F87171", "#A78BFA", "#34D399",
];

function getDatesForRange(range: string) {
  const end = new Date();
  const start = new Date();
  if (range === "30D") start.setDate(end.getDate() - 30);
  else if (range === "90D") start.setDate(end.getDate() - 90);
  else if (range === "1Y") start.setFullYear(end.getFullYear() - 1);
  else start.setFullYear(2020);
  return { start: start.toISOString(), end: end.toISOString() };
}

export default function AdminReports() {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState<"revenue" | "engagement">("revenue");
  const [dateRange, setDateRange] = useState("1Y");
  const [selectedOrganizerId, setSelectedOrganizerId] = useState<string>("all");
  const [selectedEventId, setSelectedEventId] = useState<string>("all");

  const { start, end } = getDatesForRange(dateRange);

  const { data: report, isLoading } = useQuery({
    queryKey: ["admin-revenue-report", dateRange],
    queryFn: () => adminApi.getRevenueReport(start, end),
    enabled: activeTab === "revenue",
  });

  const { data: organizersRes } = useQuery({
    queryKey: ["admin-organizers-all"],
    queryFn: () => adminApi.getOrganizers({ limit: 1000 }),
    enabled: activeTab === "engagement",
  });

  const { data: eventsRes } = useQuery({
    queryKey: ["admin-events-all"],
    queryFn: () => adminApi.getEvents({ limit: 1000 }),
    enabled: activeTab === "engagement",
  });

  const { data: engagementStats, isLoading: engagementLoading } = useQuery({
    queryKey: ["admin-reports-engagement", selectedOrganizerId, selectedEventId, dateRange],
    queryFn: () =>
      adminApi.getEngagementReportStats(
        selectedOrganizerId,
        selectedEventId,
        start,
        end
      ),
    enabled: activeTab === "engagement",
  });

  console.log("Engagement Stats:", engagementStats);

  const formatCurrency = (val: number) =>
    new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: "INR",
      maximumFractionDigits: 0,
    }).format(val);

  return (
    <AdminLayout title="Reports">
      {}
      <div className="flex border-b border-gray-200 mb-8 gap-8">
        <button
          className={`pb-4 text-sm font-semibold transition-colors relative ${
            activeTab === "revenue" ? "text-[#111827]" : "text-[#8b9098] hover:text-[#111827]"
          }`}
          onClick={() => setActiveTab("revenue")}
        >
          Revenue Report
          {activeTab === "revenue" && (
            <span className="absolute bottom-0 left-0 w-full h-[2px] bg-[#e53e5d]" />
          )}
        </button>
        <button
          className={`pb-4 text-sm font-semibold transition-colors relative ${
            activeTab === "engagement" ? "text-[#111827]" : "text-[#8b9098] hover:text-[#111827]"
          }`}
          onClick={() => setActiveTab("engagement")}
        >
          User Engagement
          {activeTab === "engagement" && (
            <span className="absolute bottom-0 left-0 w-full h-[2px] bg-[#e53e5d]" />
          )}
        </button>
      </div>

      {activeTab === "revenue" && (
        <>
          {isLoading ? (
            <div className="flex h-64 items-center justify-center">
              <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
            </div>
          ) : (
            <div className="flex flex-col gap-8">
              {}
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
                {[
                  {
                    label: "Revenue Today",
                    value: formatCurrency(report?.revenue_today?.value ?? 0),
                    pct: report?.revenue_today?.percentage ?? 0,
                    sub: "vs yesterday",
                  },
                  {
                    label: "Revenue This Month",
                    value: formatCurrency(report?.revenue_this_month?.value ?? 0),
                    pct: report?.revenue_this_month?.percentage ?? 0,
                    sub: "vs last month",
                  },
                  {
                    label: "Total Revenue",
                    value: formatCurrency(report?.total_revenue?.value ?? 0),
                    pct: report?.total_revenue?.percentage ?? 0,
                    sub: "year over year",
                  },
                  {
                    label: "Growth Rate",
                    value: `${(report?.growth_rate?.value ?? 0) >= 0 ? "+" : ""}${(report?.growth_rate?.value ?? 0).toFixed(1)}%`,
                    pct: null,
                    sub: "vs last month",
                    valueColor:
                      (report?.growth_rate?.value ?? 0) >= 0
                        ? "text-emerald-500"
                        : "text-red-500",
                  },
                ].map((card) => (
                  <div
                    key={card.label}
                    className="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex flex-col gap-2"
                  >
                    <p className="text-xs font-medium text-[#8b9098]">{card.label}</p>
                    <h3
                      className={`text-2xl font-bold ${
                        (card as any).valueColor ?? "text-[#111827]"
                      }`}
                    >
                      {card.value}
                    </h3>
                    {card.pct !== null && (
                      <div className="flex items-center gap-1.5">
                        <span
                          className={`text-xs font-semibold ${
                            card.pct >= 0 ? "text-emerald-500" : "text-red-500"
                          }`}
                        >
                          {card.pct >= 0 ? "+" : ""}{card.pct.toFixed(1)}%
                        </span>
                        <span className="text-xs text-[#8b9098]">{card.sub}</span>
                      </div>
                    )}
                    {card.pct === null && (
                      <span className="text-xs text-[#8b9098]">{card.sub}</span>
                    )}
                  </div>
                ))}
              </div>

              {}
              <div className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
                <div className="flex items-start justify-between mb-6">
                  <div>
                    <h3 className="text-lg font-bold text-[#111827]">Revenue over time</h3>
                    <p className="text-sm text-[#8b9098] mt-0.5">
                      Tracks gross revenue across selected date range
                    </p>
                  </div>
                  <div className="relative">
                    <select
                      value={dateRange}
                      onChange={(e) => setDateRange(e.target.value)}
                      className="appearance-none bg-white border border-gray-200 text-sm font-medium text-[#111827] rounded-lg px-4 py-2 pr-8 focus:outline-none cursor-pointer"
                    >
                      {DATE_RANGES.map((r) => (
                        <option key={r.value} value={r.value}>
                          {r.label}
                        </option>
                      ))}
                    </select>
                    <ChevronDown className="absolute right-2 top-1/2 -translate-y-1/2 w-4 h-4 text-[#8b9098] pointer-events-none" />
                  </div>
                </div>

                <div className="h-80">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart
                      data={report?.revenue_over_time ?? []}
                      margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
                    >
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f3f4f6" />
                      <XAxis
                        dataKey="date"
                        axisLine={false}
                        tickLine={false}
                        tick={{ fill: "#8b9098", fontSize: 12 }}
                        dy={10}
                      />
                      <YAxis
                        axisLine={false}
                        tickLine={false}
                        tick={{ fill: "#8b9098", fontSize: 12 }}
                      />
                      <RechartsTooltip
                        contentStyle={{ borderRadius: "12px", border: "none", boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1)" }}
                        formatter={(val: any) => formatCurrency(Number(val))}
                      />
                      <Line
                        type="monotone"
                        dataKey="amount"
                        stroke="#111827"
                        strokeWidth={2.5}
                        dot={{ r: 4, fill: "#111827", strokeWidth: 2, stroke: "#111827" }}
                        activeDot={{ r: 6 }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </div>

              {}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {}
                <div className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
                  <h3 className="text-lg font-bold text-[#111827]">Revenue Breakdown</h3>
                  <p className="text-sm text-[#8b9098] mt-0.5 mb-6">
                    Category-wise breakdown of the revenue generated on the platform
                  </p>

                  {(report?.category_breakdown?.length ?? 0) === 0 ? (
                    <div className="h-48 flex items-center justify-center text-sm text-[#8b9098]">
                      No category data yet.
                    </div>
                  ) : (
                    <div className="flex items-center gap-6">
                      <div className="w-[180px] h-[180px] shrink-0">
                        <ResponsiveContainer width="100%" height="100%">
                          <PieChart>
                            <Pie
                              data={report?.category_breakdown}
                              dataKey="revenue"
                              nameKey="category"
                              cx="50%"
                              cy="50%"
                              innerRadius={55}
                              outerRadius={80}
                              paddingAngle={2}
                            >
                              {report?.category_breakdown?.map((_, i) => (
                                <Cell key={i} fill={CATEGORY_COLORS[i % CATEGORY_COLORS.length]} />
                              ))}
                            </Pie>
                            <RechartsTooltip
                              formatter={(val: any) => formatCurrency(Number(val))}
                              contentStyle={{ borderRadius: "12px", border: "none", boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1)" }}
                            />
                          </PieChart>
                        </ResponsiveContainer>
                      </div>

                      <div className="flex flex-col gap-2 flex-1 min-w-0">
                        {report?.category_breakdown?.map((cat, i) => (
                          <div key={cat.category} className="flex items-center justify-between gap-2 text-sm">
                            <div className="flex items-center gap-2 min-w-0">
                              <span
                                className="w-2.5 h-2.5 rounded-full shrink-0"
                                style={{ backgroundColor: CATEGORY_COLORS[i % CATEGORY_COLORS.length] }}
                              />
                              <span className="text-[#111827] truncate">{cat.category}</span>
                            </div>
                            <span className="text-[#8b9098] font-medium shrink-0">
                              {formatCurrency(cat.revenue)}
                            </span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {}
                <div className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
                  <h3 className="text-lg font-bold text-[#111827]">Refund Analytics</h3>
                  <p className="text-sm text-[#8b9098] mt-0.5 mb-4">Monthly refund rate</p>

                  <div className="mb-4">
                    <div className="text-2xl font-bold text-[#111827]">
                      {formatCurrency(report?.refund_total?.value ?? 0)}
                    </div>
                    <div className="flex items-center gap-1.5 mt-1">
                      <span
                        className={`text-xs font-semibold ${
                          (report?.refund_total?.percentage ?? 0) >= 0
                            ? "text-emerald-500"
                            : "text-red-500"
                        }`}
                      >
                        {(report?.refund_total?.percentage ?? 0) >= 0 ? "+" : ""}
                        {(report?.refund_total?.percentage ?? 0).toFixed(1)}% this month
                      </span>
                    </div>
                  </div>

                  <div className="h-40">
                    <ResponsiveContainer width="100%" height="100%">
                      <BarChart
                        data={report?.refund_analytics ?? []}
                        margin={{ top: 5, right: 5, left: -30, bottom: 0 }}
                        barSize={24}
                      >
                        <XAxis
                          dataKey="month"
                          axisLine={false}
                          tickLine={false}
                          tick={{ fill: "#8b9098", fontSize: 11 }}
                        />
                        <YAxis hide />
                        <RechartsTooltip
                          formatter={(val: any) => formatCurrency(Number(val))}
                          contentStyle={{ borderRadius: "12px", border: "none", boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1)" }}
                        />
                        <Bar dataKey="amount" fill="#e53e5d" radius={[4, 4, 0, 0]} />
                      </BarChart>
                    </ResponsiveContainer>
                  </div>
                </div>
              </div>

              {}
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                {}
                <div className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
                  <h3 className="text-lg font-bold text-[#111827] mb-4">Top Organizers By Revenue</h3>

                  <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm">
                      <thead>
                        <tr className="text-[#8b9098] text-xs border-b border-gray-100">
                          <th className="pb-3 font-semibold">Organizer</th>
                          <th className="pb-3 font-semibold">Revenue</th>
                          <th className="pb-3 font-semibold text-center">Active Events</th>
                          <th className="pb-3 font-semibold text-center">Pending Events</th>
                          <th className="pb-3 font-semibold text-right">Avg Rating</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-50">
                        {(report?.top_organizers?.length ?? 0) === 0 ? (
                          <tr>
                            <td colSpan={5} className="py-6 text-center text-[#8b9098]">
                              No data yet.
                            </td>
                          </tr>
                        ) : (
                          report?.top_organizers?.map((org, i) => (
                            <tr key={i} className="hover:bg-[#f8fafc] transition-colors">
                              <td className="py-3 font-medium text-[#111827]">{org.name || "—"}</td>
                              <td className="py-3 text-[#111827]">{formatCurrency(org.revenue)}</td>
                              <td className="py-3 text-center text-[#111827]">{org.active_events}</td>
                              <td className="py-3 text-center text-[#111827]">{org.pending_events}</td>
                              <td className="py-3 text-right text-[#111827]">
                                {org.avg_event_rating > 0 ? org.avg_event_rating : "—"}
                              </td>
                            </tr>
                          ))
                        )}
                      </tbody>
                    </table>
                  </div>

                  <button
                    onClick={() => navigate("/admin/organizers")}
                    className="mt-4 text-sm font-semibold text-[#e53e5d] hover:underline"
                  >
                    View All Organizers &rarr;
                  </button>
                </div>

                {}
                <div className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
                  <h3 className="text-lg font-bold text-[#111827] mb-4">Top Spending Users</h3>

                  <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm">
                      <thead>
                        <tr className="text-[#8b9098] text-xs border-b border-gray-100">
                          <th className="pb-3 font-semibold">User</th>
                          <th className="pb-3 font-semibold text-center">Events Attended</th>
                          <th className="pb-3 font-semibold text-right">Total Spent</th>
                        </tr>
                      </thead>
                      <tbody className="divide-y divide-gray-50">
                        {(report?.top_users?.length ?? 0) === 0 ? (
                          <tr>
                            <td colSpan={3} className="py-6 text-center text-[#8b9098]">
                              No data yet.
                            </td>
                          </tr>
                        ) : (
                          report?.top_users?.map((user, i) => (
                            <tr key={i} className="hover:bg-[#f8fafc] transition-colors">
                              <td className="py-3 font-medium text-[#111827]">{user.name || "—"}</td>
                              <td className="py-3 text-center text-[#111827]">{user.events_attended}</td>
                              <td className="py-3 text-right text-[#111827]">{formatCurrency(user.total_spent)}</td>
                            </tr>
                          ))
                        )}
                      </tbody>
                    </table>
                  </div>

                  <button
                    onClick={() => navigate("/admin/users")}
                    className="mt-4 text-sm font-semibold text-[#e53e5d] hover:underline"
                  >
                    View All Users &rarr;
                  </button>
                </div>
              </div>
            </div>
          )}
        </>
      )}

      {activeTab === "engagement" && (
          <div className="flex flex-col gap-8">
            {}
            <div className="flex flex-wrap gap-4 items-center">
              <div className="relative">
                <select
                  value={selectedOrganizerId}
                  onChange={(e) => {
                    setSelectedOrganizerId(e.target.value);
                    setSelectedEventId("all");
                  }}
                  className="appearance-none bg-white border border-gray-200 text-[#111827] text-sm font-medium rounded-lg px-4 py-2.5 pr-10 focus:outline-none focus:ring-2 focus:ring-[#e53e5d]/20 transition-all cursor-pointer min-w-[200px]"
                >
                  <option value="all">All Organizers</option>
                  {(organizersRes as any)?.organizers?.map((org: any) => (
                    <option key={org.id} value={org.id}>
                      {org.organization_name || org.name || org.id}
                    </option>
                  ))}
                </select>
                <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#8b9098] pointer-events-none" />
              </div>

              <div className="relative">
                <select
                  value={selectedEventId}
                  onChange={(e) => setSelectedEventId(e.target.value)}
                  className="appearance-none bg-white border border-gray-200 text-[#111827] text-sm font-medium rounded-lg px-4 py-2.5 pr-10 focus:outline-none focus:ring-2 focus:ring-[#e53e5d]/20 transition-all cursor-pointer min-w-[200px]"
                >
                  <option value="all">All Events</option>
                  {(eventsRes as any)?.events?.filter((e: any) => selectedOrganizerId === "all" || e.organizer_id === selectedOrganizerId)?.map((evt: any) => (
                    <option key={evt.id} value={evt.id}>
                      {evt.title}
                    </option>
                  ))}
                </select>
                <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#8b9098] pointer-events-none" />
              </div>

              <div className="relative">
                <select
                  value={dateRange}
                  onChange={(e) => setDateRange(e.target.value)}
                  className="appearance-none bg-white border border-gray-200 text-[#111827] text-sm font-medium rounded-lg px-4 py-2.5 pr-10 focus:outline-none focus:ring-2 focus:ring-[#e53e5d]/20 transition-all cursor-pointer min-w-[200px]"
                >
                  {DATE_RANGES.map((r) => (
                    <option key={r.value} value={r.value}>{r.label}</option>
                  ))}
                </select>
                <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#8b9098] pointer-events-none" />
              </div>
            </div>

            {engagementLoading ? (
              <div className="flex h-64 items-center justify-center">
                <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
              </div>
            ) : (
              <>
                <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
                  {}
                  <div className="lg:col-span-5 flex flex-col gap-4">
                    <div className="flex gap-4">
                      {}
                      <div className="flex-1 bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
                        <p className="text-sm font-medium text-[#8b9098] mb-4">Page Views</p>
                        <h3 className="text-3xl font-bold text-[#111827] mb-2">
                          {new Intl.NumberFormat("en-IN").format(
                            (engagementStats as any)?.page_views?.value ?? 0
                          )}
                        </h3>
                        <div className="flex items-center gap-1.5 mt-auto">
                          <span className={`text-sm font-semibold ${
                            ((engagementStats as any)?.page_views?.percentage ?? 0) >= 0 ? "text-emerald-500" : "text-red-500"
                          }`}>
                            {((engagementStats as any)?.page_views?.percentage ?? 0) > 0 ? "+" : ""}
                            {((engagementStats as any)?.page_views?.percentage ?? 0).toFixed(1)}%
                          </span>
                          <span className="text-xs text-[#8b9098]">vs prior period</span>
                        </div>
                      </div>

                      {}
                      <div className="flex-1 bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
                        <p className="text-sm font-medium text-[#8b9098] mb-4">Conversion Rate</p>
                        <h3 className="text-3xl font-bold text-[#111827] mb-2">
                          {((engagementStats as any)?.conversion_rate?.value ?? 0).toFixed(1)}%
                        </h3>
                        <div className="flex items-center gap-1.5 mt-auto">
                          <span className={`text-sm font-semibold ${
                            ((engagementStats as any)?.conversion_rate?.percentage ?? 0) >= 0 ? "text-emerald-500" : "text-red-500"
                          }`}>
                            {((engagementStats as any)?.conversion_rate?.percentage ?? 0) > 0 ? "+" : ""}
                            {((engagementStats as any)?.conversion_rate?.percentage ?? 0).toFixed(1)}%
                          </span>
                          <span className="text-xs text-[#8b9098]">vs prior period</span>
                        </div>
                      </div>
                    </div>

                    <button
                      onClick={() => {
                        if (engagementStats?.user_journey) {
                          exportToCSV(engagementStats.user_journey, `admin_engagement_journey_${dateRange}`);
                        }
                      }}
                      className="w-full bg-white border border-gray-200 rounded-xl py-3.5 flex items-center justify-center gap-2 text-sm font-semibold text-[#111827] hover:bg-gray-50 transition-colors shadow-sm"
                    >
                      Export As <Download className="w-4 h-4 ml-1" />
                    </button>
                  </div>

                  {}
                  <div className="lg:col-span-7 bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
                    <h3 className="text-lg font-bold text-[#111827]">User Journey</h3>
                    <p className="text-xs text-[#8b9098] mt-1 mb-6">
                      Tracks the steps user take from visiting the event page to completing a purchase
                    </p>

                    <div className="flex flex-col gap-3">
                      {((engagementStats as any)?.user_journey ?? []).map((step: any, i: number) => {
                        const isFirst = i === 0;
                        const intensity = isFirst ? 1 : (step.percentage / 100);
                        const bgColor = isFirst
                          ? "bg-gray-100 text-gray-700"
                          : i === 1
                          ? "bg-[#f9c0cb] text-[#8b1a2e]"
                          : i === 2
                          ? "bg-[#fca5a5] text-[#7f1d1d]"
                          : i === 3
                          ? "bg-[#f87171] text-white"
                          : "bg-[#e53e5d] text-white";

                        return (
                          <div key={step.label} className="flex items-center gap-3">
                            <div className="flex flex-col items-center w-3 shrink-0">
                              <div className={`w-2.5 h-2.5 rounded-full ${
                                isFirst ? "bg-gray-300" : "bg-[#e53e5d]"
                              }`} />
                              {i < ((engagementStats as any)?.user_journey?.length ?? 0) - 1 && (
                                <div className="w-0.5 h-6 bg-[#e53e5d] mt-0.5" />
                              )}
                            </div>

                            <div
                              className={`flex-1 flex items-center justify-between px-4 py-2.5 rounded-lg text-sm font-semibold ${bgColor}`}
                              style={{
                                width: isFirst ? "100%" : `${Math.max(intensity * 100, 15)}%`,
                                minWidth: "60%",
                              }}
                            >
                              <span>{step.label}</span>
                              <span>{new Intl.NumberFormat("en-IN").format(step.count)}</span>
                            </div>

                            {!isFirst && (
                              <span className="text-sm font-semibold text-[#111827] w-12 text-right shrink-0">
                                {step.percentage.toFixed(0)}%
                              </span>
                            )}
                          </div>
                        );
                      })}

                      {((engagementStats as any)?.user_journey?.length ?? 0) === 0 && (
                        <div className="py-8 text-center text-sm text-gray-500">
                          No engagement data for this period.
                        </div>
                      )}
                    </div>
                  </div>
                </div>

                {}
                <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
                  <div className="flex items-start justify-between mb-6">
                    <div>
                      <h3 className="text-lg font-bold text-[#111827]">Peak Usage</h3>
                      <p className="text-sm text-[#8b9098] mt-1">Visualize user engagement</p>
                    </div>
                  </div>

                  <div className="h-80 w-full">
                    {((engagementStats as any)?.peak_usage?.length ?? 0) === 0 ? (
                      <div className="h-full flex items-center justify-center text-[#8b9098] text-sm">
                        No peak usage data available.
                      </div>
                    ) : (
                      <ResponsiveContainer width="100%" height="100%">
                        <BarChart
                          data={(engagementStats as any)?.peak_usage ?? []}
                          margin={{ top: 10, right: 10, left: -20, bottom: 20 }}
                          barGap={4}
                        >
                          <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f3f4f6" />
                          <XAxis
                            dataKey="label"
                            axisLine={false}
                            tickLine={false}
                            tick={{ fill: "#8b9098", fontSize: 12 }}
                            dy={10}
                            label={{ value: "Day of the week", position: "insideBottom", offset: -10, fill: "#8b9098", fontSize: 12 }}
                          />
                          <YAxis
                            axisLine={false}
                            tickLine={false}
                            tick={{ fill: "#8b9098", fontSize: 12 }}
                            label={{ value: "User Count", angle: -90, position: "insideLeft", offset: 20, fill: "#8b9098", fontSize: 12 }}
                          />
                          <RechartsTooltip
                            contentStyle={{
                              borderRadius: "12px",
                              border: "none",
                              boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1)",
                            }}
                          />
                          {}
                          <Bar dataKey="viewing" name="Views" fill="#e53e5d" radius={[4, 4, 0, 0]} />
                          <Bar dataKey="bookings" name="Bookings" fill="#10b981" radius={[4, 4, 0, 0]} />
                        </BarChart>
                      </ResponsiveContainer>
                    )}
                  </div>
                </div>
              </>
            )}
          </div>
      )}
    </AdminLayout>
  );
}
