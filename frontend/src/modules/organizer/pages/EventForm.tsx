import React, { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Upload, X, Trash2 } from "lucide-react";
import { organizerApi } from "../api";
import { adminApi } from "../../admin/api";
import type { CreateEventPayload, TicketInput, PersonnelInput } from "../api";
import OrganizerLayout from "../components/OrganizerLayout";
import { useQuery } from "@tanstack/react-query";

export default function EventForm() {
  const { eventId } = useParams();
  const isEditMode = Boolean(eventId);
  const navigate = useNavigate();

  const [title, setTitle] = useState("");
  const [category, setCategory] = useState("");
  const [startTime, setStartTime] = useState("");
  const [endTime, setEndTime] = useState("");
  const [status, setStatus] = useState("draft");
  const [tags, setTags] = useState<string[]>([]);
  const [tagInput, setTagInput] = useState("");

  const [coverImageFile, setCoverImageFile] = useState<File | null>(null);
  const [coverImagePreview, setCoverImagePreview] = useState<string>("");

  const [venueName, setVenueName] = useState("");
  const [city, setCity] = useState("");
  const [venueAddress, setVenueAddress] = useState("");
  const [mapUrl, setMapUrl] = useState("");

  const [entryType, setEntryType] = useState("Paid");
  const [freeQuantity, setFreeQuantity] = useState(0);
  const [description, setDescription] = useState("");
  const [terms, setTerms] = useState("");

  const [tickets, setTickets] = useState<TicketInput[]>([
    { name: "", price: 0, total_quantity: 0 },
  ]);

  const [personnel, setPersonnel] = useState<PersonnelInput[]>([]);

  const [submitting, setSubmitting] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");
  const [loadingInitial, setLoadingInitial] = useState(isEditMode);

  const { data: categories, isLoading: loadingCategories } = useQuery({
    queryKey: ["categories"],
    queryFn: () => adminApi.getCategories(),
  });

  useEffect(() => {
    if (!isEditMode && categories && categories.length > 0 && !category) {
      setCategory(categories[0].name);
    }
  }, [categories, isEditMode, category]);

  useEffect(() => {
    if (isEditMode && eventId) {
      const fetchEvent = async () => {
        try {
          const data = await organizerApi.getEventBySlug(eventId);
          if (data) {
            const { event, details, ticket_types, personnels } = data;

            setTitle(event.title || "");
            setCategory(event.category || "");

            if (event.start_time) {
              const st = new Date(event.start_time);
              st.setMinutes(st.getMinutes() - st.getTimezoneOffset());
              setStartTime(st.toISOString().slice(0, 16));
            }

            if (event.end_time) {
              const et = new Date(event.end_time);
              et.setMinutes(et.getMinutes() - et.getTimezoneOffset());
              setEndTime(et.toISOString().slice(0, 16));
            }

            if (event.tags) {
              setTags(event.tags.split(',').filter(Boolean));
            }

            if (event.status) {
              setStatus(event.status);
            }

            if (event.cover_image_url) {
              const url = event.cover_image_url;
              setCoverImagePreview(url.startsWith("http") ? url : `${import.meta.env.VITE_API_BASE_URL || ""}${url}`);
            }

            setVenueName(event.venue_name || "");
            setCity(event.city || "");

            if (details) {
              setVenueAddress(details.venue_address || "");
              setMapUrl(details.map_url || "");
              setDescription(details.description || "");
              setTerms(details.terms_and_conditions || "");
            }

            if (ticket_types && ticket_types.length > 0) {
              if (ticket_types.length === 1 && Number(ticket_types[0].price) === 0) {
                setEntryType("Free");
                setFreeQuantity(ticket_types[0].total_quantity);
              } else {
                setEntryType("Paid");
                setTickets(ticket_types.map((t: any) => ({
                    id: t.id,
                    name: t.name,
                    price: t.price,
                    total_quantity: t.total_quantity
                  })));
              }
            }

            if (personnels && personnels.length > 0) {
              setPersonnel(personnels.map((p: any) => ({
                id: p.id,
                name: p.name,
                role: p.role,
                image: p.image,
                profile_link: p.profile_link
              })));
            }
          }
          setLoadingInitial(false);
        } catch (err) {
          console.error(err);
          setLoadingInitial(false);
        }
      };
      fetchEvent();
    }
  }, [eventId, isEditMode]);

  const handleAddTag = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && tagInput.trim()) {
      e.preventDefault();
      if (!tags.includes(tagInput.trim())) setTags([...tags, tagInput.trim()]);
      setTagInput("");
    }
  };

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setCoverImageFile(file);
      setCoverImagePreview(URL.createObjectURL(file));
    }
  };

  const handleAddTicket = () => {
    setTickets([...tickets, { name: "", price: 0, total_quantity: 0 }]);
  };

  const handleRemoveTicket = (index: number) => {
    const newTickets = [...tickets];
    newTickets.splice(index, 1);
    setTickets(newTickets);
  };

  const updateTicket = (index: number, field: keyof TicketInput, value: string | number) => {
    const newTickets = [...tickets];
    newTickets[index] = { ...newTickets[index], [field]: value };
    setTickets(newTickets);
  };

  const totalCapacity = entryType === "Free" ? Number(freeQuantity) : tickets.reduce((acc, t) => acc + (Number(t.total_quantity) || 0), 0);
  const totalPrice = entryType === "Free" ? 0 : tickets.reduce((acc, t) => acc + ((Number(t.price) || 0) * (Number(t.total_quantity) || 0)), 0);

  const handleAddPersonnel = () => {
    setPersonnel([...personnel, { name: "", role: "", image: "", profile_link: "" }]);
  };

  const handleRemovePersonnel = (index: number) => {
    const newP = [...personnel];
    newP.splice(index, 1);
    setPersonnel(newP);
  };

  const updatePersonnel = (index: number, field: keyof PersonnelInput, value: string) => {
    const newP = [...personnel];
    newP[index] = { ...newP[index], [field]: value };
    setPersonnel(newP);
  };

  const handlePersonnelImageUpload = async (index: number, file: File) => {
    try {
        const data = await organizerApi.uploadImage(file);
        updatePersonnel(index, "image", data.url);
    } catch (err) {
        console.error("Failed to upload personnel image", err);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setErrorMsg("");

    const now = new Date();
    const stDate = new Date(startTime);
    if (stDate < now) {
      setErrorMsg("Event start time cannot be in the past");
      setSubmitting(false);
      return;
    }

    try {
      let finalCoverUrl = "";
      if (coverImageFile) {
        const data = await organizerApi.uploadImage(coverImageFile);
        finalCoverUrl = data.url;
      }

      const payload = {
        title,
        city,
        venue_name: venueName,
        category,
        start_time: new Date(startTime).toISOString(),
        end_time: endTime ? new Date(endTime).toISOString() : new Date(new Date(startTime).getTime() + 2 * 60 * 60 * 1000).toISOString(),
        tags: tags,
        cover_image_url: finalCoverUrl || (coverImagePreview
          ? (coverImagePreview.startsWith("http")
              ? coverImagePreview
              : coverImagePreview.replace(import.meta.env.VITE_API_BASE_URL || "", ""))
          : ""),
        status: (isEditMode && status.toLowerCase() === "approved") ? "draft" : status,
        details: {
          description,
          venue_address: venueAddress,
          map_url: mapUrl,
          total_capacity: totalCapacity,
          terms_and_conditions: terms,
        },
        ticket_types: entryType === "Free"
          ? [{ name: "Free Entry", price: 0, total_quantity: Number(freeQuantity) }]
          : tickets.filter(t => t.name.trim() !== "").map(t => ({
              ...t,
              price: Number(t.price),
              total_quantity: Number(t.total_quantity)
            })),
        key_personnel: personnel,
      };

      if (isEditMode && eventId) {
        await organizerApi.updateEvent(eventId, payload);
      } else {
        await organizerApi.createEvent(payload as CreateEventPayload);
      }

      navigate("/organizer/events", {
        state: { toastMessage: isEditMode ? "Event successfully updated!" : "Event successfully created!" }
      });

    } catch (err: any) {
      console.error(err);
      setErrorMsg(err.response?.data?.message || err.message || "Failed to submit event");
    } finally {
      setSubmitting(false);
    }
  };

  if (loadingInitial || loadingCategories) {
    return (
        <div className="min-h-screen bg-[#f8f9fa] flex items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
    );
  }

  return (
    <OrganizerLayout activeTab="My Events">
      <div className="py-10 px-10 max-w-4xl mx-auto">
        <div className="flex flex-col gap-4 mb-8">
          <button type="button" onClick={() => navigate(-1)} className="flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-gray-900 transition w-fit">
            <span aria-hidden="true">&larr;</span> Back
          </button>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{isEditMode ? "Update the event" : "Create a new event for the community"}</h1>
            <p className="text-sm text-gray-500 mt-1">Events are published only post approval.</p>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="bg-white rounded-[24px] p-10 shadow-sm border border-gray-100 min-h-[600px] flex flex-col gap-8">

          {errorMsg && (
            <div className="p-4 bg-red-50 text-red-600 rounded-xl text-sm font-medium border border-red-100">
               {errorMsg}
            </div>
          )}

          
          <div>
              <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Event Title</label>
              <input
                type="text"
                required
                placeholder="e.g. Premium Roy by Shreya"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 focus:ring-1 focus:ring-gray-400"
              />
          </div>

          <div className="grid grid-cols-2 gap-6">
              <div>
                  <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Category</label>
                  <select
                    value={category}
                    onChange={(e) => setCategory(e.target.value)}
                    className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400"
                  >
                      {categories?.map((c) => (
                        <option key={c.id} value={c.name}>{c.name}</option>
                      ))}
                      {!categories?.length && <option value="Comedy">Comedy</option>}
                  </select>
              </div>
              <div className="flex gap-4 col-span-2">
                  <div className="flex-1">
                      <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Start Date & Time</label>
                      <input
                        type="datetime-local"
                        required
                        min={new Date().toISOString().slice(0, 16)}
                        value={startTime}
                        onChange={(e) => setStartTime(e.target.value)}
                        className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400"
                      />
                  </div>
                  <div className="flex-1">
                      <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">End Date & Time</label>
                      <input
                        type="datetime-local"
                        min={startTime || new Date().toISOString().slice(0, 16)}
                        value={endTime}
                        onChange={(e) => setEndTime(e.target.value)}
                        className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400"
                      />
                  </div>
              </div>
          </div>

          <div>
              <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Tags</label>
              <div className="flex flex-wrap items-center gap-2 border border-gray-200 rounded-xl p-2 min-h-[50px]">
                  {tags.map((tag, idx) => (
                      <div key={idx} className="bg-gray-100 text-gray-700 px-3 py-1.5 rounded-lg text-sm flex items-center gap-2">
                          {tag}
                          <button type="button" onClick={() => setTags(tags.filter((_, i) => i !== idx))}><X className="w-3 h-3 text-red-500"/></button>
                      </div>
                  ))}
                  <input
                    type="text"
                    placeholder={tags.length === 0 ? "Add tags and press enter..." : "..."}
                    value={tagInput}
                    onChange={(e) => setTagInput(e.target.value)}
                    onKeyDown={handleAddTag}
                    className="flex-1 min-w-[120px] outline-none text-sm px-2 py-1 bg-transparent"
                  />
              </div>
          </div>

          <div>
              <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Event Poster</label>
              <div className="border-2 border-dashed border-gray-200 rounded-2xl p-8 flex flex-col items-center justify-center relative bg-gray-50 group hover:border-[#e53e5d]/30 transition-colors">
                  {coverImagePreview ? (
                      <div className="relative w-full aspect-video rounded-xl overflow-hidden shadow-sm">
                          <img src={coverImagePreview} alt="Cover Preview" className="w-full h-full object-cover" />
                          <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                              <span className="text-white text-sm font-medium">Click to change format</span>
                          </div>
                      </div>
                  ) : (
                      <>
                          <div className="w-12 h-12 bg-white rounded-xl shadow-sm border border-gray-100 flex items-center justify-center mb-3">
                              <Upload className="w-5 h-5 text-gray-400" />
                          </div>
                          <p className="text-sm font-medium text-gray-900"><span className="text-[#e53e5d]">Upload Image</span></p>
                          <p className="text-xs text-gray-500 mt-1">Supported formats: JPG, JPEG, PNG, WEBP</p>
                          <p className="text-xs text-gray-400">Maximum file size: 10 MB</p>
                      </>
                  )}
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleImageUpload}
                    className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                  />
              </div>
          </div>

          <div className="pt-6 border-t border-gray-100">
              <label className="block text-sm font-bold text-gray-900 tracking-wide mb-6">VENUE DETAILS</label>
              <div className="grid grid-cols-2 gap-6 mb-6">
                <div>
                     <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Venue Name</label>
                     <input type="text" required placeholder="JLN Stadium" value={venueName} onChange={e => setVenueName(e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400" />
                </div>
                <div>
                     <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase text-left">Venue Location (City)</label>
                     <input
                       type="text"
                       required
                       placeholder="e.g. Kochi"
                       value={city}
                       onChange={e => setCity(e.target.value)}
                       className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400"
                     />
                </div>
              </div>
              <div className="mb-6">
                  <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Full Address</label>
                  <input type="text" required placeholder="Stadium Road, Kaloor..." value={venueAddress} onChange={e => setVenueAddress(e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400" />
              </div>
              <div>
                  <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Map URL (Opt.)</label>
                  <input type="url" placeholder="https://maps.google.com/..." value={mapUrl} onChange={e => setMapUrl(e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400" />
              </div>
          </div>

          <div className="pt-6 border-t border-gray-100 flex gap-6 items-center">
              <label className="block text-sm font-bold text-gray-900 tracking-wide">ENTRY TYPE</label>
              <div className="flex gap-4">
                  <label className="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
                      <input type="radio" name="entry_type" checked={entryType === "Free"} onChange={() => setEntryType("Free")} className="accent-[#e53e5d] w-4 h-4" /> Free Entry
                  </label>
                  <label className="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
                      <input type="radio" name="entry_type" checked={entryType === "Paid"} onChange={() => setEntryType("Paid")} className="accent-[#e53e5d] w-4 h-4" /> Paid Entry
                  </label>
              </div>
          </div>

          {entryType === "Free" && (
              <div className="pt-6 border-t border-gray-100">
                  <div className="bg-gray-50 border border-gray-200 rounded-2xl p-6">
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 items-center">
                          <div>
                              <label className="block text-sm font-bold text-gray-900 tracking-wide mb-1 uppercase">Total Capacity</label>
                              <p className="text-xs text-gray-500 mb-4">Set the maximum number of people who can attend for free.</p>
                              <input
                                type="number"
                                required
                                min="1"
                                placeholder="e.g. 500"
                                value={freeQuantity || ''}
                                onChange={e => setFreeQuantity(Number(e.target.value))}
                                className="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-900 bg-white outline-none focus:border-gray-400"
                              />
                          </div>
                      </div>
                  </div>
              </div>
          )}

          {entryType === "Paid" && (
              <div className="pt-6 border-t border-gray-100">
                  <div className="flex justify-between items-center mb-6">
                      <label className="block text-sm font-bold text-gray-900 tracking-wide uppercase">Ticket Configuration</label>
                      <button type="button" onClick={handleAddTicket} className="text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-lg transition-colors border border-gray-200">
                          + Add Ticket Type
                      </button>
                  </div>

                  <div className="flex flex-col gap-4">
                      {tickets.map((ticket, idx) => (
                          <div key={idx} className="bg-gray-50 border border-gray-200 rounded-2xl p-6 relative group">
                              <div className="flex justify-between items-center mb-4">
                                  <h4 className="text-xs font-bold text-gray-400 uppercase tracking-widest">Ticket Type {idx + 1}</h4>
                                  {tickets.length > 1 && (
                                      <button type="button" onClick={() => handleRemoveTicket(idx)} className="text-red-400 hover:text-red-600 transition-colors p-1 opacity-0 group-hover:opacity-100 focus:opacity-100">
                                          <Trash2 className="w-4 h-4" />
                                      </button>
                                  )}
                              </div>
                              <div className="grid grid-cols-1 gap-4 mb-4">
                                  <div>
                                      <label className="block text-[11px] font-bold text-gray-500 tracking-wider mb-1.5 uppercase">Ticket Name</label>
                                      <input type="text" required placeholder="Premium" value={ticket.name} onChange={e => updateTicket(idx, "name", e.target.value)} disabled={isEditMode && status !== "draft" && !!(ticket as any).id} className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white outline-none disabled:bg-gray-100 disabled:text-gray-500 focus:border-gray-400" />
                                  </div>
                              </div>
                               <div className="grid grid-cols-2 gap-4">
                                   <div>
                                       <div className="flex justify-between items-center mb-1.5">
                                           <label className="block text-[11px] font-bold text-gray-500 tracking-wider uppercase">Price (₹)</label>
                                       </div>
                                       <input type="number" required min="1" placeholder="5000" value={ticket.price || ''} onChange={e => updateTicket(idx, "price", e.target.value)} disabled={isEditMode && status !== "draft" && !!(ticket as any).id} className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white outline-none disabled:bg-gray-100 disabled:text-gray-500 focus:border-gray-400" />
                                   </div>
                                   <div>
                                       <label className="block text-[11px] font-bold text-gray-500 tracking-wider mb-1.5 uppercase">Quantity</label>
                                       <input type="number" required min="1" placeholder="100" value={ticket.total_quantity || ''} onChange={e => updateTicket(idx, "total_quantity", e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white outline-none focus:border-gray-400" />
                                   </div>
                               </div>
                           </div>
                       ))}
                   </div>

                   <div className="grid grid-cols-2 gap-6 mt-6">
                       <div>
                           <label className="block text-xs font-bold text-gray-500 tracking-wider mb-2 uppercase">Total Capacity</label>
                           <div className="w-full bg-gray-100 border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-700 font-semibold">{totalCapacity}</div>
                           <p className="text-[10px] text-gray-400 mt-2">Automatically calculated based on the ticket categories provided.</p>
                       </div>
                       <div>
                           <label className="block text-xs font-bold text-gray-500 tracking-wider mb-2 uppercase">Total Price Value</label>
                           <div className="w-full bg-gray-100 border border-gray-200 rounded-xl px-4 py-3 text-sm text-gray-700 font-semibold">
                               ₹ {totalPrice.toLocaleString()}
                           </div>
                       </div>
                   </div>
              </div>
          )}

          <div className="pt-6 border-t border-gray-100">
              <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">About Event</label>
              <textarea required rows={5} placeholder="Give information regarding the event..." value={description} onChange={e => setDescription(e.target.value)} className="w-full border border-gray-200 rounded-xl p-4 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 resize-y" />
          </div>

          <div>
              <label className="block text-xs font-semibold text-gray-500 tracking-wider mb-2 uppercase">Terms & Conditions</label>
              <textarea rows={4} placeholder="Add terms and conditions for your attendees." value={terms} onChange={e => setTerms(e.target.value)} className="w-full border border-gray-200 rounded-xl p-4 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 resize-y" />
          </div>

          <div className="pt-6 border-t border-gray-100">
              <div className="flex justify-between items-center mb-6">
                  <label className="block text-sm font-bold text-gray-900 tracking-wide uppercase">Key Personnel</label>
                  <button type="button" onClick={handleAddPersonnel} className="text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-lg transition-colors border border-gray-200">
                      + Add Personnel
                  </button>
              </div>

              <div className="flex flex-col gap-6">
                  {personnel.map((p, idx) => (
                      <div key={idx} className="bg-white border border-gray-200 rounded-2xl p-6 relative group">
                           <div className="flex justify-between items-center mb-4">
                                  <h4 className="text-xs font-bold text-gray-400 uppercase tracking-widest">Personnel {idx + 1}</h4>
                                  <button type="button" onClick={() => handleRemovePersonnel(idx)} className="text-red-400 hover:text-red-600 transition-colors p-1 opacity-0 group-hover:opacity-100 focus:opacity-100">
                                      <Trash2 className="w-4 h-4" />
                                  </button>
                           </div>
                           <div className="grid grid-cols-2 gap-6 mb-6">
                                <div>
                                    <label className="block text-[11px] font-bold text-gray-500 tracking-wider mb-1.5 uppercase">Name</label>
                                    <input type="text" required placeholder="Joe Smith" value={p.name} onChange={e => updatePersonnel(idx, "name", e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white outline-none focus:border-gray-400" />
                                </div>
                                <div>
                                    <label className="block text-[11px] font-bold text-gray-500 tracking-wider mb-1.5 uppercase">Role</label>
                                    <input type="text" required placeholder="Lead Performer" value={p.role} onChange={e => updatePersonnel(idx, "role", e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white outline-none focus:border-gray-400" />
                                </div>
                           </div>

                           <div className="mb-6">
                                <label className="block text-[11px] font-bold text-gray-500 tracking-wider mb-1.5 uppercase">Thumbnail Image</label>
                                <div className="border border-dashed border-gray-300 rounded-xl p-4 flex flex-col items-center justify-center relative bg-gray-50 hover:bg-gray-100 transition-colors cursor-pointer group/upload">
                                    {p.image ? (
                                         <div className="w-16 h-16 rounded-full overflow-hidden shadow-sm border border-gray-200 relative">
                                            <img src={p.image.startsWith("http") ? p.image : `${import.meta.env.VITE_API_BASE_URL}${p.image}`} alt={p.name} className="w-full h-full object-cover"/>
                                         </div>
                                    ) : (
                                        <>
                                            <div className="w-10 h-10 bg-white rounded-full shadow-sm border border-gray-200 flex items-center justify-center mb-2">
                                                <Upload className="w-4 h-4 text-gray-400" />
                                            </div>
                                            <p className="text-xs font-semibold text-[#e53e5d] group-hover/upload:text-[#d03550]">Upload Image</p>
                                        </>
                                    )}
                                    <input type="file" onChange={(e) => { if(e.target.files?.[0]) handlePersonnelImageUpload(idx, e.target.files[0]) }} className="absolute inset-0 w-full h-full opacity-0 cursor-pointer" />
                                </div>
                           </div>

                             <div>
                                <label className="block text-[11px] font-bold text-gray-500 tracking-wider mb-1.5 uppercase">Link to Profile (Opt.)</label>
                                <input type="url" placeholder="https://instagram.com/..." value={p.profile_link || ''} onChange={e => updatePersonnel(idx, "profile_link", e.target.value)} className="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm bg-white outline-none focus:border-gray-400" />
                            </div>
                      </div>
                  ))}
                  {personnel.length === 0 && (
                      <p className="text-sm text-gray-400 text-center py-4 border border-dashed border-gray-200 rounded-xl">No personnel added yet.</p>
                  )}
              </div>
          </div>

          <div className="mt-6 pt-4">
              <button disabled={submitting} type="submit" className="w-full bg-[#0b101e] hover:bg-black text-white py-4 rounded-xl text-sm font-semibold tracking-wide flex items-center justify-center transition-colors disabled:opacity-70 disabled:cursor-not-allowed">
                  {submitting ? (
                      <div className="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin" />
                  ) : (
                      isEditMode ? "Save Changes" : "Submit Event"
                  )}
              </button>
              <p className="text-[11px] text-gray-400 text-center mt-4 px-10 leading-relaxed">
                  By submitting, you agree to the EVNTX <span className="font-semibold text-gray-600">Terms & Conditions</span> for event creation and management.
              </p>
          </div>

        </form>
      </div>
    </OrganizerLayout>
  );
}
