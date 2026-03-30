import { useEffect, useRef, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Bell } from "lucide-react";
import { notificationsApi } from "../../modules/notifications/api";
import { useNotifications } from "../../modules/notifications/hooks";
import { useAuthStore } from "../../modules/auth/store/authStore";

type NotificationMenuProps = {
  buttonClassName?: string;
  iconClassName?: string;
  panelClassName?: string;
};

export default function NotificationMenu({
  buttonClassName,
  iconClassName,
  panelClassName,
}: NotificationMenuProps) {
  const { isAuthenticated } = useAuthStore();
  const [open, setOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const queryClient = useQueryClient();
  const { data, isLoading } = useNotifications(1, 6, isAuthenticated);

  const markAsReadMutation = useMutation({
    mutationFn: (notificationId: string) => notificationsApi.markAsRead(notificationId),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });

  const markAllAsReadMutation = useMutation({
    mutationFn: () => notificationsApi.markAllAsRead(),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });

  const clearAllMutation = useMutation({
    mutationFn: () => notificationsApi.clearAll(),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["notifications"] });
    },
  });

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  if (!isAuthenticated) {
    return null;
  }

  const unreadCount = data?.unread_count ?? 0;
  const notifications = data?.notifications ?? [];

  return (
    <div className="relative" ref={menuRef}>
      <button
        type="button"
        onClick={() => setOpen((current) => !current)}
        className={buttonClassName ?? "relative rounded-full p-2 text-[#6c7480] transition hover:bg-[#f5f5f5]"}
      >
        <Bell className={iconClassName ?? "h-5 w-5"} />
        {unreadCount > 0 ? (
          <span className="absolute -right-1 -top-1 min-w-[20px] rounded-full bg-[#ff445d] px-1.5 py-0.5 text-center text-[10px] font-semibold leading-none text-white">
            {unreadCount > 99 ? "99+" : unreadCount}
          </span>
        ) : null}
      </button>

      {open ? (
        <div className={panelClassName ?? "absolute right-0 top-12 w-80 rounded-2xl border border-gray-100 bg-white p-3 shadow-[0_16px_40px_rgba(15,23,42,0.12)]"}>
          <div className="flex items-center justify-between px-2 py-1">
            <div className="text-sm font-semibold text-[#111827]">Notifications</div>
            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => clearAllMutation.mutate()}
                disabled={clearAllMutation.isPending || notifications.length === 0}
                className="text-xs font-medium text-[#6b7280] hover:text-[#111827] disabled:cursor-not-allowed disabled:text-gray-300 transition"
              >
                Clear all
              </button>
              <button
                type="button"
                onClick={() => markAllAsReadMutation.mutate()}
                disabled={markAllAsReadMutation.isPending || unreadCount === 0}
                className="text-xs font-medium text-[#ff445d] disabled:cursor-not-allowed disabled:text-gray-300"
              >
                Mark all read
              </button>
            </div>
          </div>

          <div className="mt-2 flex max-h-96 flex-col gap-2 overflow-y-auto">
            {isLoading ? (
              <div className="rounded-xl bg-[#f8fafc] px-3 py-4 text-sm text-[#6b7280]">
                Loading notifications...
              </div>
            ) : notifications.length === 0 ? (
              <div className="rounded-xl bg-[#f8fafc] px-3 py-4 text-sm text-[#6b7280]">
                No notifications yet.
              </div>
            ) : (
              notifications.slice(0, 5).map((notification) => (
                <button
                  key={notification.id}
                  type="button"
                  onClick={() => {
                    if (!notification.is_read) {
                      markAsReadMutation.mutate(notification.id);
                    }
                  }}
                  className={`rounded-xl border px-3 py-3 text-left transition ${
                    notification.is_read
                      ? "border-transparent bg-[#f8fafc]"
                      : "border-[#ffd4db] bg-[#fff7f8]"
                  }`}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="text-sm font-medium text-[#111827]">{notification.title}</div>
                      <div className="mt-1 text-xs leading-5 text-[#6b7280]">{notification.message}</div>
                    </div>
                    {!notification.is_read ? (
                      <span className="mt-1 h-2.5 w-2.5 shrink-0 rounded-full bg-[#ff445d]" />
                    ) : null}
                  </div>
                  <div className="mt-2 text-[11px] uppercase tracking-[0.14em] text-[#9ca3af]">
                    {new Date(notification.created_at).toLocaleString("en-IN", {
                      day: "2-digit",
                      month: "short",
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </div>
                </button>
              ))
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}
