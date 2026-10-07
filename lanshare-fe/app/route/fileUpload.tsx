import {Form, useActionData} from "react-router";
import {useState} from "react";
import type {Route} from "../../.react-router/types/app/route/+types/fileUpload";

export async function clientAction({request}: Route.ClientActionArgs) {
    const formData = await request.formData();

    const text = String(formData.get("text") ?? "").trim();
    const url = String(formData.get("url") ?? "").trim();
    const file = formData.get("file");
    const expiresAt = String(formData.get("expires_at") ?? "");

    const hasFile = file instanceof File && file.size > 0;

    if (!text && !url && !hasFile) {
        return {
            error: "Ievadiet tekstu, izvēlieties failu vai ievadiet URL.",
        };
    }

    if (url) {
        try {
            const parsedUrl = new URL(url);

            if (
                parsedUrl.protocol !== "http:" &&
                parsedUrl.protocol !== "https:"
            ) {
                return {
                    error: "URL jābūt HTTP vai HTTPS adresei.",
                };
            }
        } catch {
            return {
                error: "Ievadiet derīgu URL.",
            };
        }
    }

    if (!expiresAt) {
        return {
            error: "Norādiet derīguma termiņu.",
        };
    }

    const expiration = new Date(expiresAt);

    if (Number.isNaN(expiration.getTime())) {
        return {
            error: "Norādīts nederīgs derīguma termiņš.",
        };
    }

    if (expiration <= new Date()) {
        return {
            error: "Derīguma termiņam jābūt nākotnē.",
        };
    }

    formData.set("text", text);
    formData.set("url", url);
    formData.set("expires_at", expiration.toISOString());

    if (!hasFile) {
        formData.delete("file");
    }

    const response = await fetch("/api/files/upload", {
        method: "POST",
        body: formData,
        credentials: "include",
    });

    return await response.json();
}

export default function FileUpload() {
    const actionData = useActionData<typeof clientAction>();
    const [expiresAt, setExpiresAt] = useState(() => getExpiration(24));

    const minExpiresAt = formatDateTimeLocal(new Date());

    return (
        <Form
            method="post"
            encType="multipart/form-data"
            className="mx-auto flex w-full max-w-2xl flex-col gap-6"
        >
            <div className="flex flex-col gap-2">
                <label
                    htmlFor="text"
                    className="text-sm font-medium text-gray-900"
                >
                    Teksts
                </label>

                <textarea
                    id="text"
                    name="text"
                    rows={5}
                    maxLength={10000}
                    placeholder="Ievadiet tekstu..."
                    className="w-full resize-y rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition placeholder:text-gray-400 focus:border-gray-500 focus:ring-2 focus:ring-gray-200"
                />
            </div>

            <div className="flex flex-col gap-2">
                <label
                    htmlFor="file"
                    className="text-sm font-medium text-gray-900"
                >
                    Fails
                </label>

                <input
                    id="file"
                    name="file"
                    type="file"
                    className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-700 file:mr-4 file:rounded-md file:border-0 file:bg-gray-100 file:px-3 file:py-2 file:text-sm file:font-medium file:text-gray-700 hover:file:bg-gray-200"
                />
            </div>

            <div className="flex flex-col gap-2">
                <label
                    htmlFor="url"
                    className="text-sm font-medium text-gray-900"
                >
                    URL
                </label>

                <input
                    id="url"
                    name="url"
                    type="url"
                    placeholder="https://example.com"
                    className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition placeholder:text-gray-400 focus:border-gray-500 focus:ring-2 focus:ring-gray-200"
                />
            </div>

            <div className="flex flex-col gap-3">
                <label
                    htmlFor="expires_at"
                    className="text-sm font-medium text-gray-900"
                >
                    Derīguma termiņš
                </label>

                <div className="flex flex-wrap gap-2">
                    <button
                        type="button"
                        onClick={() => setExpiresAt(getExpiration(1))}
                        className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
                    >
                        1 stunda
                    </button>

                    <button
                        type="button"
                        onClick={() => setExpiresAt(getExpiration(24))}
                        className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
                    >
                        1 diena
                    </button>

                    <button
                        type="button"
                        onClick={() => setExpiresAt(getExpiration(72))}
                        className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
                    >
                        3 dienas
                    </button>

                    <button
                        type="button"
                        onClick={() => setExpiresAt(getExpiration(168))}
                        className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
                    >
                        7 dienas
                    </button>
                </div>

                <input
                    id="expires_at"
                    name="expires_at"
                    type="datetime-local"
                    value={expiresAt}
                    min={minExpiresAt}
                    required
                    onChange={(event) => setExpiresAt(event.target.value)}
                    className="w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 outline-none transition focus:border-gray-500 focus:ring-2 focus:ring-gray-200"
                />
            </div>

            {actionData?.error && (
                <p className="text-sm text-red-600">
                    {actionData.error}
                </p>
            )}

            <div className="flex justify-end">
                <button
                    type="submit"
                    className="rounded-lg bg-gray-900 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-gray-800 active:bg-gray-950"
                >
                    Augšupielādēt
                </button>
            </div>
        </Form>
    );
}

function formatDateTimeLocal(date: Date) {
    const offset = date.getTimezoneOffset();
    const localDate = new Date(date.getTime() - offset * 60 * 1000);

    return localDate.toISOString().slice(0, 16);
}

function getExpiration(hours: number) {
    return formatDateTimeLocal(
        new Date(Date.now() + hours * 60 * 60 * 1000),
    );
}