import {Link, useLoaderData} from "react-router";
import {useState} from "react";
import type  {Route} from "../../.react-router/types/app/route/+types/fileView";
import type {file} from "~/types/file";
export async function clientLoader({params}: Route.ClientLoaderArgs) {
    const response = await fetch(`/api/files/${params.id}`);
    return await response.json();
}
interface LoaderData {
    file?: file,
    error?: string,
}
export default function FileView() {
    const loaderData = useLoaderData<LoaderData>();
    const file = loaderData.file;
    const downloadUrl = file?.filename? "/api/files/" + file.id + "/download" : null;
    const [copied, setCopied] = useState(false);

    const copyLink = async () => {
        if(!file) return
        await navigator.clipboard.writeText(
            `${window.location.origin}/share/${file.id}`,
        );

        setCopied(true);
        setTimeout(() => setCopied(false), 1500);
    };

    return (
        <main className="mx-auto w-full max-w-3xl px-4 py-8">
            {loaderData.error && (
                <div className="flex min-h-100 items-center justify-center">
                    <div className="w-full max-w-lg rounded-2xl border border-gray-300 bg-gray-100/80 p-8 text-center shadow-lg backdrop-blur-sm">
                        <div className="mx-auto mb-5 flex h-16 w-16 items-center justify-center rounded-full bg-red-100 text-3xl">
                            404
                        </div>

                        <h1 className="text-2xl font-bold text-gray-900">
                            Fails netika atrasts
                        </h1>

                        <p className="mt-2 text-sm leading-6 text-gray-600">
                            Šis fails, iespējams, ir dzēsts vai arī tā derīguma termiņš ir beidzies.
                        </p>

                        <Link
                            to="/"
                            className="mt-6 inline-flex rounded-lg bg-gray-900 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-gray-800"
                        >
                            Apskatīt citus failus
                        </Link>
                    </div>
                </div>
            )}
            {file && (
                <div className="rounded-2xl border border-gray-300 bg-gray-100/80 shadow-lg backdrop-blur-sm">
                <div className="flex items-center justify-between gap-4 border-b border-gray-300 px-6 py-5">
                    <div className="min-w-0">
                        <h1 className="truncate text-2xl font-bold text-gray-900">
                            {file.filename ?? "Dalītais fails"}
                        </h1>

                        <p className="mt-1 text-sm text-gray-500">
                            Augšupielādēts {formatDate(file?.uploaded_at)}
                        </p>
                    </div>

                    <div className="relative shrink-0">
                        <button
                            type="button"
                            onClick={copyLink}
                            className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium text-gray-700 transition hover:bg-gray-50"
                        >
                            Kopēt saiti
                        </button>

                        {copied && (
                            <div className="absolute bottom-full left-1/2 mb-2 -translate-x-1/2 whitespace-nowrap rounded-md bg-gray-900 px-2 py-1 text-xs font-medium text-white shadow-md">
                                Nokopēts
                            </div>
                        )}
                    </div>
                </div>

                <div className="p-6">
                    {file.text && (
                        <div className="rounded-xl border my-3 border-gray-300 bg-white p-5">
                            <p className="whitespace-pre-wrap wrap-break-word text-sm leading-6 text-gray-800">
                                {file.text}
                            </p>
                        </div>
                    )}

                    {file.url && (
                        <div className="rounded-xl border my-3 border-gray-300 bg-white p-5">
                            <p className="mb-2 text-xs font-medium uppercase tracking-wide text-gray-400">
                                URL
                            </p>

                            <a
                                href={file.url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="break-all text-blue-600 hover:underline"
                            >
                                {file.url}
                            </a>
                        </div>
                    )}

                    {downloadUrl && (
                        <div className="flex flex-col my-3 items-center rounded-xl border border-gray-300 bg-white px-6 py-10 text-center">
                            <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-gray-100 text-2xl">
                                📄
                            </div>

                            <h2 className="max-w-full truncate text-lg font-semibold text-gray-900">
                                {file.filename ?? "Fails"}
                            </h2>

                            {file.size !== null && (
                                <p className="mt-1 text-sm text-gray-500">
                                    {formatFileSize(file!.size!)}
                                </p>
                            )}

                            <a
                                href={downloadUrl}
                                className="mt-5 rounded-lg bg-gray-900 px-5 py-2.5 text-sm font-medium text-white transition hover:bg-gray-800"
                            >
                                Lejupielādēt
                            </a>
                        </div>
                    )}

                    {!file.text && !file.url && !downloadUrl && (
                        <div className="rounded-xl border border-gray-300 bg-white px-6 py-10 text-center">
                            <p className="text-sm text-gray-500">
                                Šis fails nav pieejams.
                            </p>
                        </div>
                    )}
                </div>

                <div className="flex items-center justify-between border-t border-gray-300 px-6 py-4">
                    <div className="text-xs text-gray-500">
                        Derīgs līdz {formatDate(file.expires_at!)}
                    </div>

                    <Link
                        to="/"
                        className="text-sm font-medium text-gray-700 hover:text-gray-900"
                    >
                        Atpakaļ
                    </Link>
                </div>
            </div>
                )}
        </main>
    );
}

function formatFileSize(bytes: number) {
    if (bytes === 0) return "0 B";

    const units = ["B", "KB", "MB", "GB"];
    const index = Math.floor(Math.log(bytes) / Math.log(1024));

    return `${(bytes / Math.pow(1024, index)).toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

function formatDate(date: string) {
    return new Date(date).toLocaleString("lv-LV");
}