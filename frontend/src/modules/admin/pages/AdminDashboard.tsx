import { useState } from "react";
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
} from "recharts";
import {
  TrendingUp,
  Users,
  UserCheck,
  CalendarDays,
  RotateCcw,
  UserPlus,
  ClipboardList,
  Zap,
  Loader2,
} from "lucide-react";

interface StatCardProps {
  title: string;
  value: string;
  percentage: number;
  subtitle: string;
  icon: React.ElementType;
  iconBg: string;
  iconColor: string;
}

function StatCard({ title, value, percentage, subtitle, icon: Icon, iconBg, iconColor }: StatCardProps) {
  const isPositive = percentage >= 0;
  const isNeutral = percentage === 0;

  return (
    <div className="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${iconBg}`}>
          <Icon className={`w-5 h-5 ${iconColor}`} />
        </div>
        <span className="text-sm text-[#8b9098] font-medium">{title}</span>
      </div>
      <div>
        <div className="text-2xl font-bold text-[#111827] leading-tight">{value}</div>
        <div className="flex items-center gap-1.5 mt-1.5">
          {!isNeutral && (
            <span className={`text-xs font-semibold ${isPositive ? "text-emerald-500" : "text-red-500"}`}>
              {isPositive ? "+" : ""}{percentage.toFixed(1)}%
            </span>
          )}
          <span className="text-xs text-[#8b9098]">{subtitle}</span>
        </div>
      </div>
    </div>
  );
}

export default function AdminDashboard() {
  const [revenueSpan, setRevenueSpan] = useState("1Y");

  const { data: stats, isLoading } = useQuery({
    queryKey: ["admin-dashboard-stats", revenueSpan],
    queryFn: () => adminApi.getDashboardStats({ span: revenueSpan }),
  });

  const formatCurrency = (val: number) =>
    new Intl.NumberFormat("en-IN", {
      style: "currency",
      currency: "INR",
      maximumFractionDigits: 0,
    }).format(val);

  const formatNumber = (val: number) =>
    new Intl.NumberFormat("en-IN").format(Math.round(val));

  const statCards = stats
    ? [
        {
          title: "Revenue",
          value: formatCurrency(stats.revenue.value),
          percentage: stats.revenue.percentage,
          subtitle: "this month",
          icon: TrendingUp,
          iconBg: "bg-indigo-50",
          iconColor: "text-indigo-600",
        },
        {
          title: "Total Users",
          value: formatNumber(stats.total_users.value),
          percentage: stats.total_users.percentage,
          subtitle: "this month",
          icon: Users,
          iconBg: "bg-blue-50",
          iconColor: "text-blue-500",
        },
        {
          title: "Total Organizers",
          value: formatNumber(stats.total_organizers.value),
          percentage: stats.total_organizers.percentage,
          subtitle: "this month",
          icon: UserCheck,
          iconBg: "bg-teal-50",
          iconColor: "text-teal-500",
        },
        {
          title: "Total Events",
          value: formatNumber(stats.total_events.value),
          percentage: stats.total_events.percentage,
          subtitle: "this month",
          icon: CalendarDays,
          iconBg: "bg-sky-50",
          iconColor: "text-sky-500",
        },
        {
          title: "Refund Rate",
          value: `${stats.refund_rate.value.toFixed(1)}%`,
          percentage: stats.refund_rate.percentage,
          subtitle: stats.refund_rate.subtitle || "this month",
          icon: RotateCcw,
          iconBg: "bg-rose-50",
          iconColor: "text-rose-500",
        },
        {
          title: "User Growth",
          value: formatNumber(stats.user_growth.value),
          percentage: stats.user_growth.percentage,
          subtitle: "this month",
          icon: UserPlus,
          iconBg: "bg-emerald-50",
          iconColor: "text-emerald-500",
        },
        {
          title: "Pending Approvals",
          value: formatNumber(stats.pending_approvals.value),
          percentage: 0,
          subtitle: "events awaiting review",
          icon: ClipboardList,
          iconBg: "bg-amber-50",
          iconColor: "text-amber-500",
        },
        {
          title: "Active Events",
          value: formatNumber(stats.active_events.value),
          percentage: stats.active_events.percentage,
          subtitle: "this month",
          icon: Zap,
          iconBg: "bg-violet-50",
          iconColor: "text-violet-500",
        },
      ]
    : [];


  return (
    <AdminLayout title="Dashboard">
      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
          <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
        </div>
      ) : (
        <div className="flex flex-col gap-8">
          
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
            {statCards.map((card) => (
              <StatCard key={card.title} {...card} />
            ))}
          </div>

          
          <div className="bg-white rounded-2xl p-6 border border-gray-100 shadow-sm">
            <div className="mb-6">
              <h3 className="text-lg font-bold text-[#111827]">Performance Summary</h3>
              <p className="text-sm text-[#8b9098] mt-0.5">
                Revenue increased{" "}
                {stats?.revenue?.percentage && stats.revenue.percentage > 0
                  ? `${stats.revenue.percentage.toFixed(0)}%`
                  : "—"}{" "}
                this month.
              </p>
            </div>

            
            <div className="flex items-center justify-between mb-4">
              <h4 className="font-semibold text-[#111827]">Revenue Overview</h4>
              <div className="flex items-center gap-2">
                <div className="flex bg-gray-100 rounded-lg p-1">
                  {["7D", "30D", "90D", "1Y", "ALL"].map((span) => (
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
            </div>

            <div className="h-80 w-full">
              {(!stats?.revenue_overview || stats.revenue_overview.length === 0) ? (
                <div className="h-full flex items-center justify-center text-[#8b9098] text-sm">
                  No revenue data available.
                </div>
              ) : (
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart
                    data={stats.revenue_overview}
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
                      dot={{ r: 4, fill: "#111827", strokeWidth: 2, stroke: "#111827" }}
                      activeDot={{ r: 6 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              )}
            </div>
          </div>
        </div>
      )}
    </AdminLayout>
  );
}
