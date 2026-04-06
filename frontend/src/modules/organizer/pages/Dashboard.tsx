import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import OrganizerLayout from "../components/OrganizerLayout";
import { organizerApi } from "../api";
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
  Cell
} from "recharts";
import { Banknote, Ticket, CalendarCheck, Clock, Loader2 } from "lucide-react";

export default function Dashboard() {
  const navigate = useNavigate();
  const [revenueSpan, setRevenueSpan] = useState("1Y");

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ["organizer-dashboard-stats"],
    queryFn: () => organizerApi.getDashboardStats(),
  });

  const { data: events, isLoading: eventsLoading } = useQuery({
    queryKey: ["organizer-events"],
    queryFn: () => organizerApi.getOrganizerEvents(),
  });

  const formatCurrency = (val: number) =>
    new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: "INR",
      maximumFractionDigits: 0,
    }).format(val);

  const StatCard = ({
    title,
    value,
    percentage,
    subtitle,
    icon: Icon,
    colorClass,
  }: any) => {
    const isPositive = percentage > 0;
    const isNeutral = percentage === 0;

    return (
      <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
        <div className="flex justify-between items-start mb-4">
          <div className={`p-3 rounded-xl ${colorClass}`}>
            <Icon className="w-6 h-6" />
          </div>
          <p className="text-sm font-medium text-[#8b9098]">{title}</p>
        </div>
        <div className="mt-auto">
          <h3 className="text-3xl font-bold text-[#111827] mb-2">{value}</h3>
          <div className="flex items-center gap-2">
            {!isNeutral && (
              <span
                className={`text-sm font-semibold ${
                  isPositive ? "text-emerald-500" : "text-red-500"
                }`}
              >
                {isPositive ? "+" : ""}{percentage}%
              </span>
            )}
            <span className="text-xs text-[#8b9098]">{subtitle}</span>
          </div>
        </div>
      </div>
    );
  };

  const COLORS = ["#2DD4BF", "#4F46E5", "#FBBF24", "#60A5FA", "#F87171"];

  const renderStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case "pending":
        return (
          <span className="px-3 py-1 rounded-full text-xs font-semibold bg-[#FBBF24] text-white">
            Pending
          </span>
        );
      case "live":
        return (
          <span className="px-3 py-1 rounded-full text-xs font-semibold bg-white border border-[#e53e5d] text-[#e53e5d]">
            Live
          </span>
        );
      case "approved":
        return (
          <span className="px-3 py-1 rounded-full text-xs font-semibold bg-[#06b6d4] text-white">
            Approved
          </span>
        );
      default:
        return (
          <span className="px-3 py-1 rounded-full text-xs font-semibold bg-gray-100 text-gray-800 capitalize">
            {status}
          </span>
        );
    }
  };

  const isLoading = statsLoading || eventsLoading;

  return (
    <OrganizerLayout activeTab="Dashboard">
      <div className="p-8 pb-24">
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-[#111827]">Dashboard</h1>
        </div>

        {isLoading ? (
          <div className="flex h-64 items-center justify-center">
            <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
          </div>
        ) : (
          <>
            {/* STATS STRIP */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              <StatCard
                title="Total Revenue"
                value={formatCurrency(stats?.total_revenue?.value || 0)}
                percentage={stats?.total_revenue?.percentage || 0}
                subtitle="this month"
                icon={Banknote}
                colorClass="bg-indigo-50 text-indigo-600"
              />
              <StatCard
                title="Tickets Sold"
                value={new Intl.NumberFormat("en-IN").format(
                  stats?.tickets_sold?.value || 0
                )}
                percentage={stats?.tickets_sold?.percentage || 0}
                subtitle="this month"
                icon={Ticket}
                colorClass="bg-amber-50 text-amber-500"
              />
              <StatCard
                title="Active Events"
                value={stats?.active_events?.value || 0}
                percentage={stats?.active_events?.percentage || 0}
                subtitle="this month"
                icon={CalendarCheck}
                colorClass="bg-teal-50 text-teal-500"
              />
              <StatCard
                title="Pending Events"
                value={stats?.pending_events?.value || 0}
                percentage={stats?.pending_events?.percentage || 0}
                subtitle="in pipeline"
                icon={Clock}
                colorClass="bg-red-50 text-red-500"
              />
            </div>

            {/* CHARTS GRID */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
              <div className="lg:col-span-2 bg-white rounded-2xl p-6 shadow-sm border border-gray-100">
                <div className="flex justify-between items-start mb-6">
                  <div>
                    <h3 className="text-lg font-bold text-[#111827]">
                      Performance Summary
                    </h3>
                    <p className="text-sm text-[#8b9098] mt-1">
                      Revenue overview over time.
                    </p>
                  </div>
                  <div className="flex bg-gray-100 rounded-lg p-1">
                    {["7D", "30D", "90D", "1Y"].map((span) => (
                      <button
                        key={span}
                        onClick={() => setRevenueSpan(span)}
                        className={`px-3 py-1.5 text-xs font-medium rounded-md transition-all ${
                          revenueSpan === span
                            ? "bg-white text-gray-900 shadow-sm"
                            : "text-[#8b9098] hover:text-gray-900"
                        }`}
                      >
                        {span}
                      </button>
                    ))}
                  </div>
                </div>

                <div className="h-72 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <LineChart
                      data={stats?.revenue_overview || []}
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
                        stroke="#e53e5d"
                        strokeWidth={3}
                        dot={{ r: 4, fill: "#e53e5d", strokeWidth: 2, stroke: "#fff" }}
                        activeDot={{ r: 6 }}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </div>

              <div className="bg-white rounded-2xl p-6 shadow-sm border border-gray-100 flex flex-col">
                <h3 className="text-lg font-bold text-[#111827]">
                  Ticket Sales Breakdown
                </h3>
                <p className="text-xs text-[#8b9098] mt-1 mb-6">
                  Breakdown of ticket sales for active events
                </p>

                {(stats?.sales_breakdown?.length || 0) === 0 ? (
                  <div className="flex-1 flex items-center justify-center text-sm text-[#8b9098]">
                    No sales data yet
                  </div>
                ) : (
                  <>
                    <div className="h-[200px] w-full mb-6">
                      <ResponsiveContainer width="100%" height="100%">
                        <PieChart>
                          <Pie
                            data={stats?.sales_breakdown || []}
                            cx="50%"
                            cy="50%"
                            innerRadius={60}
                            outerRadius={80}
                            paddingAngle={2}
                            dataKey="value"
                          >
                            {stats?.sales_breakdown?.map((_: any, index: number) => (
                              <Cell
                                key={`cell-${index}`}
                                fill={COLORS[index % COLORS.length]}
                              />
                            ))}
                          </Pie>
                          <RechartsTooltip
                            formatter={(val: any) => formatCurrency(Number(val))}
                            contentStyle={{
                              borderRadius: "12px",
                              border: "none",
                              boxShadow: "0 10px 15px -3px rgb(0 0 0 / 0.1)",
                            }}
                          />
                        </PieChart>
                      </ResponsiveContainer>
                    </div>

                    <div className="flex flex-col gap-3 mt-auto mb-2 text-sm">
                      {stats?.sales_breakdown?.slice(0, 5).map((item: any, i: number) => (
                        <div key={i} className="flex items-center justify-between">
                          <div className="flex items-center gap-2 truncate pr-4">
                            <span
                              className="w-3 h-3 rounded-full shrink-0"
                              style={{ backgroundColor: COLORS[i % COLORS.length] }}
                            />
                            <span className="text-[#111827] truncate">
                              {item.name}
                            </span>
                          </div>
                          <span className="text-[#8b9098] font-medium shrink-0">
                            {formatCurrency(item.value)}
                          </span>
                        </div>
                      ))}
                    </div>
                  </>
                )}
              </div>
            </div>

            {/* TABLE */}
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
              <div className="overflow-x-auto">
                <table className="w-full text-left text-sm text-[#111827]">
                  <thead className="bg-[#f8fafc] text-[#8b9098]">
                    <tr>
                      <th className="px-6 py-4 font-semibold">Event Name</th>
                      <th className="px-6 py-4 font-semibold">Date</th>
                      <th className="px-6 py-4 font-semibold">Tickets Sold</th>
                      <th className="px-6 py-4 font-semibold">Revenue</th>
                      <th className="px-6 py-4 font-semibold">Status</th>
                      <th className="px-6 py-4 font-semibold">Action</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100">
                    {events?.data?.length === 0 ? (
                      <tr>
                        <td colSpan={6} className="px-6 py-8 text-center text-[#8b9098]">
                          No events found
                        </td>
                      </tr>
                    ) : (
                      events?.data?.map((event: any) => (
                        <tr
                          key={event.id}
                          className="hover:bg-[#f8fafc] transition-colors"
                        >
                          <td className="px-6 py-4 font-medium max-w-[200px] truncate">
                            {event.title}
                          </td>
                          <td className="px-6 py-4">
                            {new Date(event.start_time).toLocaleDateString("en-IN", {
                              day: "numeric",
                              month: "short",
                              year: "numeric",
                            })}
                          </td>
                          <td className="px-6 py-4">{event.tickets_sold || 0}</td>
                          <td className="px-6 py-4 text-gray-500">
                            {formatCurrency(event.revenue || 0)}
                          </td>
                          <td className="px-6 py-4">
                            {renderStatusBadge(event.status)}
                          </td>
                          <td className="px-6 py-4">
                            <button
                              onClick={() => navigate(`/organizer/events/${event.id}/edit`)}
                              className="text-sm font-semibold text-[#111827] hover:underline"
                            >
                              View
                            </button>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
              
              <div className="border-t border-gray-100 px-6 py-4 flex items-center justify-end gap-2 text-sm text-[#8b9098]">
                {/* Visual Pagination Mock matching Figma */}
                <button className="px-2 py-1 hover:text-gray-900">&lt; Prev</button>
                <button className="w-8 h-8 rounded bg-gray-200 text-gray-900 font-medium flex items-center justify-center">1</button>
                <button className="w-8 h-8 rounded hover:bg-gray-100 flex items-center justify-center">2</button>
                <button className="w-8 h-8 rounded hover:bg-gray-100 flex items-center justify-center">3</button>
                <button className="px-2 py-1 hover:text-gray-900">Next &gt;</button>
              </div>
            </div>
          </>
        )}
      </div>
    </OrganizerLayout>
  );
}
