import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import OrganizerLayout from "../components/OrganizerLayout";
import { organizerApi } from "../api";
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  Legend,
} from "recharts";
import { Download, ChevronDown, Loader2 } from "lucide-react";
import { exportToCSV } from "../../../shared/utils/csv";

export default function Reports() {
  const [activeTab, setActiveTab] = useState<"sales" | "engagement">("sales");
  const [selectedEventId, setSelectedEventId] = useState<string>("all");
  
  
  const [dateRange, setDateRange] = useState<"30D" | "90D" | "1Y" | "ALL">("30D");

  const getDatesForRange = (range: string) => {
    const end = new Date();
    const start = new Date();
    if (range === "30D") start.setDate(end.getDate() - 30);
    else if (range === "90D") start.setDate(end.getDate() - 90);
    else if (range === "1Y") start.setFullYear(end.getFullYear() - 1);
    else if (range === "ALL") start.setFullYear(2020); 
    
    return {
      start: start.toISOString(),
      end: end.toISOString()
    }
  }

  const { start, end } = getDatesForRange(dateRange);

  const { data: eventsRes } = useQuery({
    queryKey: ["organizer-events"],
    queryFn: () => organizerApi.getOrganizerEvents(),
  });

  const { data: reportStats, isLoading } = useQuery({
    queryKey: ["organizer-reports-sales", selectedEventId, dateRange],
    queryFn: () =>
      organizerApi.getSalesReportStats(
        selectedEventId,
        start,
        end
      ),
    enabled: activeTab === "sales",
  });

  const { data: engagementStats, isLoading: engagementLoading } = useQuery({
    queryKey: ["organizer-reports-engagement", selectedEventId, dateRange],
    queryFn: () =>
      organizerApi.getEngagementReportStats(
        selectedEventId,
        start,
        end
      ),
    enabled: activeTab === "engagement",
  });

  const formatCurrency = (val: number) =>
    new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: "INR",
      maximumFractionDigits: 0,
    }).format(val);

  return (
    <OrganizerLayout activeTab="Reports">
      <div className="p-8 pb-24">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-[#111827]">
            Access all your reports
          </h1>
        </div>

        {}
        <div className="flex border-b border-gray-200 mb-8 w-fit gap-8 px-2">
          <button
            className={`pb-4 text-sm font-semibold transition-colors relative ${
              activeTab === "sales"
                ? "text-[#111827]"
                : "text-[#8b9098] hover:text-[#111827]"
            }`}
            onClick={() => setActiveTab("sales")}
          >
            Sales Report
            {activeTab === "sales" && (
              <span className="absolute bottom-0 left-0 w-full h-[2px] bg-[#e53e5d]" />
            )}
          </button>
          <button
            className={`pb-4 text-sm font-semibold transition-colors relative ${
              activeTab === "engagement"
                ? "text-[#111827]"
                : "text-[#8b9098] hover:text-[#111827]"
            }`}
            onClick={() => setActiveTab("engagement")}
          >
            User Engagement
            {activeTab === "engagement" && (
              <span className="absolute bottom-0 left-0 w-full h-[2px] bg-[#e53e5d]" />
            )}
          </button>
        </div>

        {activeTab === "sales" && (
          <div className="flex flex-col gap-8">
            {}
            <div className="flex flex-wrap gap-4 items-center">
              <div className="relative">
                <select
                  value={selectedEventId}
                  onChange={(e) => setSelectedEventId(e.target.value)}
                  className="appearance-none bg-white border border-gray-200 text-[#111827] text-sm font-medium rounded-lg px-4 py-2.5 pr-10 focus:outline-none focus:ring-2 focus:ring-[#e53e5d]/20 transition-all cursor-pointer min-w-[200px]"
                >
                  <option value="all">All Events</option>
                  {eventsRes?.events?.map((evt: any) => (
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
                  onChange={(e) => setDateRange(e.target.value as any)}
                  className="appearance-none bg-white border border-gray-200 text-[#111827] text-sm font-medium rounded-lg px-4 py-2.5 pr-10 focus:outline-none focus:ring-2 focus:ring-[#e53e5d]/20 transition-all cursor-pointer min-w-[200px]"
                >
                  <option value="30D">Last 30 Days</option>
                  <option value="90D">Last 90 Days</option>
                  <option value="1Y">Last Year</option>
                  <option value="ALL">All Time</option>
                </select>
                <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#8b9098] pointer-events-none" />
              </div>
            </div>

            {isLoading ? (
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
                      <p className="text-sm font-medium text-[#8b9098] mb-4">
                        Total Revenue
                      </p>
                      <h3 className="text-3xl font-bold text-[#111827] mb-2 truncate">
                        {formatCurrency(reportStats?.total_revenue?.value || 0)}
                      </h3>
                      <div className="flex items-center gap-1.5 mt-auto">
                        <span
                          className={`text-sm font-semibold ${
                            (reportStats?.total_revenue?.percentage || 0) >= 0
                              ? "text-emerald-500"
                              : "text-red-500"
                          }`}
                        >
                          {(reportStats?.total_revenue?.percentage || 0) > 0 ? "+" : ""}
                          {(reportStats?.total_revenue?.percentage || 0).toFixed(1)}%
                        </span>
                        <span className="text-xs text-[#8b9098]">
                          vs prior period
                        </span>
                      </div>
                    </div>

                    {}
                    <div className="flex-1 bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
                      <p className="text-sm font-medium text-[#8b9098] mb-4">
                        Tickets Sold
                      </p>
                      <h3 className="text-3xl font-bold text-[#111827] mb-2">
                        {new Intl.NumberFormat("en-IN").format(
                          reportStats?.tickets_sold?.value || 0
                        )}
                      </h3>
                      <div className="flex items-center gap-1.5 mt-auto">
                        <span
                          className={`text-sm font-semibold ${
                            (reportStats?.tickets_sold?.percentage || 0) >= 0
                              ? "text-emerald-500"
                              : "text-red-500"
                          }`}
                        >
                          {(reportStats?.tickets_sold?.percentage || 0) > 0 ? "+" : ""}
                          {(reportStats?.tickets_sold?.percentage || 0).toFixed(1)}%
                        </span>
                        <span className="text-xs text-[#8b9098]">
                          vs prior period
                        </span>
                      </div>
                    </div>
                  </div>

                  <button
                    onClick={() => {
                      if (reportStats?.revenue_over_time) {
                        exportToCSV(reportStats.revenue_over_time, `organizer_revenue_${dateRange}`, [
                          { header: "Date", key: "date" },
                          { header: "Amount", key: "amount" },
                        ]);
                      }
                    }}
                    className="w-full bg-white border border-gray-200 rounded-xl py-3.5 flex items-center justify-center gap-2 text-sm font-semibold text-[#111827] hover:bg-gray-50 transition-colors shadow-sm"
                  >
                    Export As <Download className="w-4 h-4 ml-1" />
                  </button>
                </div>

                {}
                <div className="lg:col-span-7 bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
                  <h3 className="text-lg font-bold text-[#111827] mb-6">
                    Tickets sold per event
                  </h3>
                  
                  <div className="grid grid-cols-[1fr_2fr_1fr_1fr] text-xs font-semibold text-[#8b9098] pb-3 border-b border-gray-100">
                    <div className="col-span-2">Event</div>
                    <div className="text-center">Tickets Sold</div>
                    <div className="text-right">% Total</div>
                  </div>

                  <div className="flex flex-col pt-2 py-1 max-h-[180px] overflow-y-auto pr-2 custom-scrollbar">
                    {reportStats?.tickets_per_event?.length === 0 ? (
                      <div className="py-8 text-center text-sm text-gray-500">
                        No tickets sold in this period.
                      </div>
                    ) : (
                      reportStats?.tickets_per_event?.map((evt: any, i: number) => (
                        <div
                          key={i}
                          className="grid grid-cols-[1fr_2fr_1fr_1fr] items-center py-3 border-b border-gray-50 last:border-0"
                        >
                          <div className="text-sm font-medium text-[#111827] truncate pr-4">
                            {evt.name}
                          </div>
                          
                          {}
                          <div className="px-4">
                            <div className="h-3 w-full bg-gray-100 rounded-full overflow-hidden flex">
                              <div
                                className="h-full bg-[#e53e5d] rounded-full"
                                style={{ width: `${evt.percentage_total}%` }}
                              />
                            </div>
                          </div>
                          
                          <div className="text-sm text-[#111827] text-center">
                            {evt.tickets_sold}
                          </div>
                          <div className="text-sm text-[#8b9098] text-right font-medium">
                            {evt.percentage_total.toFixed(1)}%
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                </div>
              </div>

              {}
              <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 mt-8">
                <div className="mb-6">
                  <h3 className="text-lg font-bold text-[#111827]">
                    Revenue over time
                  </h3>
                  <p className="text-sm text-[#8b9098] mt-1">
                    Tracks gross revenue across selected date range
                  </p>
                </div>

                <div className="h-80 w-full">
                  {(reportStats?.revenue_over_time?.length || 0) === 0 ? (
                    <div className="h-full flex items-center justify-center text-[#8b9098] text-sm">
                      No data to chart.
                    </div>
                  ) : (
                    <ResponsiveContainer width="100%" height="100%">
                      <LineChart
                        data={reportStats?.revenue_over_time || []}
                        margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
                      >
                        <CartesianGrid
                          strokeDasharray="3 3"
                          vertical={false}
                          stroke="#f3f4f6"
                        />
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
                          contentStyle={{
                            borderRadius: "12px",
                            border: "none",
                            boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1)",
                          }}
                          formatter={(val: any) => formatCurrency(Number(val))}
                        />
                        <Line
                          type="monotone"
                          dataKey="amount"
                          stroke="#111827"
                          strokeWidth={2.5}
                          dot={{
                            r: 4,
                            fill: "#111827",
                            strokeWidth: 2,
                            stroke: "#111827",
                          }}
                          activeDot={{ r: 6 }}
                        />
                      </LineChart>
                    </ResponsiveContainer>
                  )}
                </div>
              </div>
              </>
            )}
          </div>
        )}
        
        {activeTab === "engagement" && (
          <div className="flex flex-col gap-8">
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
                            engagementStats?.page_views?.value ?? 0
                          )}
                        </h3>
                        <div className="flex items-center gap-1.5 mt-auto">
                          <span className={`text-sm font-semibold ${
                            (engagementStats?.page_views?.percentage ?? 0) >= 0 ? "text-emerald-500" : "text-red-500"
                          }`}>
                            {(engagementStats?.page_views?.percentage ?? 0) > 0 ? "+" : ""}
                            {(engagementStats?.page_views?.percentage ?? 0).toFixed(1)}%
                          </span>
                          <span className="text-xs text-[#8b9098]">vs prior period</span>
                        </div>
                      </div>

                      {}
                      <div className="flex-1 bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
                        <p className="text-sm font-medium text-[#8b9098] mb-4">Conversion Rate</p>
                        <h3 className="text-3xl font-bold text-[#111827] mb-2">
                          {(engagementStats?.conversion_rate?.value ?? 0).toFixed(1)}%
                        </h3>
                        <div className="flex items-center gap-1.5 mt-auto">
                          <span className={`text-sm font-semibold ${
                            (engagementStats?.conversion_rate?.percentage ?? 0) >= 0 ? "text-emerald-500" : "text-red-500"
                          }`}>
                            {(engagementStats?.conversion_rate?.percentage ?? 0) > 0 ? "+" : ""}
                            {(engagementStats?.conversion_rate?.percentage ?? 0).toFixed(1)}%
                          </span>
                          <span className="text-xs text-[#8b9098]">vs prior period</span>
                        </div>
                      </div>
                    </div>

                    <button
                      onClick={() => {
                        if (engagementStats?.user_journey) {
                          exportToCSV(engagementStats.user_journey, `organizer_engagement_journey_${dateRange}`, [
                            { header: "Stage", key: "label" },
                            { header: "Count", key: "count" },
                            { header: "Percentage", key: "percentage" },
                          ]);
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
                      {(engagementStats?.user_journey ?? []).map((step, i) => {
                        const isFirst = i === 0;
                        const intensity = isFirst ? 1 : (step.percentage / 100);
                        const bgColor = isFirst
                          ? "bg-gray-100 text-gray-700"
                          : i === 1
                          ? "bg-[#f9c0cb] text-[#8b1a2e]"
                          : i === 2
                          ? "bg-[#f87171] text-white"
                          : "bg-[#e53e5d] text-white";

                        return (
                          <div key={step.label} className="flex items-center gap-3">
                            {}
                            <div className="flex flex-col items-center w-3 shrink-0">
                              <div className={`w-2.5 h-2.5 rounded-full ${
                                isFirst ? "bg-gray-300" : "bg-[#e53e5d]"
                              }`} />
                              {i < (engagementStats?.user_journey?.length ?? 0) - 1 && (
                                <div className="w-0.5 h-6 bg-[#e53e5d] mt-0.5" />
                              )}
                            </div>

                            {}
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

                            {}
                            {!isFirst && (
                              <span className="text-sm font-semibold text-[#111827] w-12 text-right shrink-0">
                                {step.percentage.toFixed(0)}%
                              </span>
                            )}
                          </div>
                        );
                      })}

                      {(engagementStats?.user_journey?.length ?? 0) === 0 && (
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
                      <p className="text-sm text-[#8b9098] mt-1">Visualize user engagement on events</p>
                    </div>
                  </div>

                  <div className="h-80 w-full">
                    {(engagementStats?.peak_usage?.length ?? 0) === 0 ? (
                      <div className="h-full flex items-center justify-center text-[#8b9098] text-sm">
                        No peak usage data available.
                      </div>
                    ) : (
                      <ResponsiveContainer width="100%" height="100%">
                        <BarChart
                          data={engagementStats?.peak_usage ?? []}
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
                          <Legend
                            verticalAlign="top"
                            align="right"
                            iconType="circle"
                            wrapperStyle={{ paddingBottom: "12px", fontSize: "12px" }}
                          />
                          <Bar dataKey="checkout" name="Checkout" fill="#6b7280" radius={[4, 4, 0, 0]} />
                          <Bar dataKey="viewing" name="Viewing" fill="#e53e5d" radius={[4, 4, 0, 0]} />
                        </BarChart>
                      </ResponsiveContainer>
                    )}
                  </div>
                </div>
              </>
            )}
          </div>
        )}
      </div>
    </OrganizerLayout>
  );
}
