import { useEffect, useId, useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { Camera, CheckCircle2, Loader2, QrCode, Ticket, XCircle } from "lucide-react";
import { useParams, useSearchParams } from "react-router-dom";
import OrganizerLayout from "../components/OrganizerLayout";
import { organizerApi, type CheckInResponse } from "../api";

export default function OrganizerCheckInPage() {
  const { eventId = "" } = useParams();
  const [searchParams] = useSearchParams();
  const scannerId = useId().replace(/:/g, "-");
  const [ticketCode, setTicketCode] = useState("");
  const [scannerOpen, setScannerOpen] = useState(false);
  const [scannerError, setScannerError] = useState("");
  const [result, setResult] = useState<CheckInResponse | null>(null);
  const [errorMessage, setErrorMessage] = useState("");
  const hasScannedRef = useRef(false);

  useEffect(() => {
    const code = searchParams.get("code");
    if (code) {
      setTicketCode(code);
    }
  }, [searchParams]);

  const checkInMutation = useMutation({
    mutationFn: (code: string) => organizerApi.checkInTicket(eventId, code),
    onSuccess: (data) => {
      setResult(data);
      setErrorMessage("");
      setScannerError("");
    },
    onError: (error: any) => {
      setResult(null);
      setErrorMessage(error?.response?.data?.error?.message || "Failed to validate ticket.");
    },
  });

  const submitCheckIn = async (codeOverride?: string) => {
    const normalizedCode = (codeOverride ?? ticketCode).trim();
    if (!normalizedCode) {
      setErrorMessage("Ticket code is required.");
      setResult(null);
      return;
    }

    setTicketCode(normalizedCode);
    await checkInMutation.mutateAsync(normalizedCode);
  };

  useEffect(() => {
    if (!scannerOpen) {
      return;
    }

    let isMounted = true;
    let html5QrCode: any = null;
    hasScannedRef.current = false;
    setScannerError("");

    const startScanner = async () => {
      try {
        const { Html5Qrcode } = await import("html5-qrcode");
        if (!isMounted) {
          return;
        }

        const cameras = await Html5Qrcode.getCameras();
        if (!cameras.length) {
          setScannerError("No camera found on this device.");
          return;
        }

        html5QrCode = new Html5Qrcode(scannerId);
        await html5QrCode.start(
          cameras[0].id,
          {
            fps: 10,
            qrbox: { width: 220, height: 220 },
            aspectRatio: 1,
          },
          (decodedText: string) => {
            if (hasScannedRef.current) {
              return;
            }

            hasScannedRef.current = true;
            const normalizedCode = decodedText.trim();
            setTicketCode(normalizedCode);
            setScannerOpen(false);
            void submitCheckIn(normalizedCode);
          },
          () => undefined,
        );
      } catch {
        if (isMounted) {
          setScannerError("Unable to start QR scanner. Check camera permission and try again.");
        }
      }
    };

    void startScanner()

    return () => {
      isMounted = false;
      if (html5QrCode) {
        void html5QrCode.stop()
          .catch(() => undefined)
          .then(() => html5QrCode.clear().catch(() => undefined));
      }
    };
  }, [scannerId, scannerOpen]);

  return (
    <OrganizerLayout activeTab="My Events">
      <div className="mx-auto max-w-4xl px-4 py-10 lg:px-10">
        <div className="rounded-[28px] border border-gray-100 bg-white p-8 shadow-sm">
          <div className="flex flex-col gap-3 border-b border-gray-100 pb-6">
            <div className="flex items-center gap-3 text-[#111827]">
              <div className="rounded-2xl bg-[#eef6ff] p-3 text-[#0f62fe]">
                <Ticket className="h-6 w-6" />
              </div>
              <div>
                <h1 className="text-2xl font-semibold">Ticket Check-In</h1>
                <p className="mt-1 text-sm text-gray-500">
                  Validate tickets for this event with manual entry or QR scan.
                </p>
              </div>
            </div>
            <p className="text-xs uppercase tracking-[0.24em] text-gray-400">Event ID: {eventId}</p>
          </div>

          <div className="mt-8 grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
            <div className="rounded-[24px] border border-gray-100 bg-[#fafbfc] p-6">
              <div className="flex items-center gap-2 text-sm font-medium text-gray-700">
                <QrCode className="h-4 w-4" />
                Manual Entry
              </div>
              <div className="mt-4 flex flex-col gap-4">
                <label className="text-sm text-gray-600" htmlFor="ticket-code">
                  Ticket code
                </label>
                <input
                  id="ticket-code"
                  value={ticketCode}
                  onChange={(event) => setTicketCode(event.target.value)}
                  placeholder="Enter or scan ticket code"
                  className="w-full rounded-2xl border border-gray-200 bg-white px-4 py-3 text-sm text-gray-900 outline-none transition focus:border-[#111827]"
                />
                <div className="flex flex-col gap-3 sm:flex-row">
                  <button
                    type="button"
                    onClick={() => void submitCheckIn()}
                    disabled={checkInMutation.isPending}
                    className="inline-flex items-center justify-center gap-2 rounded-2xl bg-[#111827] px-5 py-3 text-sm font-semibold text-white transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    {checkInMutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <CheckCircle2 className="h-4 w-4" />}
                    Check In
                  </button>
                  <button
                    type="button"
                    onClick={() => setScannerOpen((current) => !current)}
                    className="inline-flex items-center justify-center gap-2 rounded-2xl border border-gray-200 bg-white px-5 py-3 text-sm font-semibold text-gray-800 transition hover:bg-gray-50"
                  >
                    <Camera className="h-4 w-4" />
                    {scannerOpen ? "Close Scanner" : "Open QR Scanner"}
                  </button>
                </div>
              </div>

              {scannerOpen ? (
                <div className="mt-6 rounded-[24px] border border-dashed border-gray-300 bg-white p-4">
                  <div id={scannerId} className="overflow-hidden rounded-2xl" />
                  <p className="mt-3 text-xs text-gray-500">
                    Point the camera at the ticket QR code to auto-submit the check-in request.
                  </p>
                </div>
              ) : null}

              {scannerError ? (
                <div className="mt-4 rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
                  {scannerError}
                </div>
              ) : null}
            </div>

            <div className="rounded-[24px] border border-gray-100 bg-white p-6">
              <div className="text-sm font-medium text-gray-700">Validation Result</div>

              {result ? (
                <div className="mt-5 rounded-[24px] border border-emerald-200 bg-emerald-50 p-5">
                  <div className="flex items-center gap-2 text-emerald-700">
                    <CheckCircle2 className="h-5 w-5" />
                    <span className="text-base font-semibold">Ticket Validated</span>
                  </div>
                  <div className="mt-4 space-y-3 text-sm text-emerald-900">
                    <div>
                      <div className="text-xs uppercase tracking-[0.2em] text-emerald-600">Ticket code</div>
                      <div className="mt-1 font-medium">{result.ticket_code}</div>
                    </div>
                    <div>
                      <div className="text-xs uppercase tracking-[0.2em] text-emerald-600">Status</div>
                      <div className="mt-1 font-medium capitalize">{result.status}</div>
                    </div>
                    <div>
                      <div className="text-xs uppercase tracking-[0.2em] text-emerald-600">Checked in at</div>
                      <div className="mt-1 font-medium">
                        {new Date(result.checked_in_at).toLocaleString("en-IN", {
                          day: "2-digit",
                          month: "short",
                          year: "numeric",
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </div>
                    </div>
                  </div>
                </div>
              ) : errorMessage ? (
                <div className="mt-5 rounded-[24px] border border-rose-200 bg-rose-50 p-5">
                  <div className="flex items-center gap-2 text-rose-700">
                    <XCircle className="h-5 w-5" />
                    <span className="text-base font-semibold">Check-in Failed</span>
                  </div>
                  <p className="mt-3 text-sm text-rose-900">{errorMessage}</p>
                </div>
              ) : (
                <div className="mt-5 rounded-[24px] border border-gray-100 bg-[#f8fafc] p-5 text-sm text-gray-500">
                  Validate a ticket to see the latest check-in result here.
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </OrganizerLayout>
  );
}
