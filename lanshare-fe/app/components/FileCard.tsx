import {useState} from "react";
import type {file} from "~/types/file";
import {useFetcher} from "react-router";

export default function FileCard({
                                     file,
                                     type,
                                     route,
                                 }: {
    file: file;
    type: "owned" | "public";
    route: "/files" | "/share";
}) {
    const [copied, setCopied] = useState(false);

    let fileName = file.filename;


    if (!fileName) {
        fileName = file.text?.slice(0, 50);
    }

    if (!fileName) {
        fileName = file.url;
    }
    const fetcher = useFetcher();


    const handleDelete = () => {
        fetcher.submit(null, {
            method: "delete",
            action: `/files/delete/${file.id}`,
        });
    };
    const copyShareLink = async () => {
        await navigator.clipboard.writeText(
            `${window.location.origin}/share/${file.id}`,
        );

        setCopied(true);

        setTimeout(() => {
            setCopied(false);
        }, 1500);
    };

    return (
        <article
            className="group rounded-xl border border-gray-300/80 bg-gray-100/90 p-4 shadow-sm backdrop-blur-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-gray-400 hover:bg-white hover:shadow-md">
            <div className="flex items-center justify-between gap-4">
                <a href={`${route}/${file.id}`} className="min-w-0 flex-1">
                    <h3 className="truncate font-semibold text-gray-900">
                        {fileName}
                    </h3>

                    <div className="mt-1 flex gap-3 text-xs text-gray-500">
                        {file.size !== undefined && (
                            <span>{formatFileSize(file.size)}</span>
                        )}

                        <span>
                            {new Date(file.uploaded_at).toLocaleDateString("lv-LV")}
                        </span>
                    </div>
                </a>

                <div
                    className="flex shrink-0 items-center gap-2 opacity-0 transition-opacity duration-200 group-hover:opacity-100">
                    <a
                        href={`/files/${file.id}`}
                        className="rounded-lg bg-gray-800 px-3 py-2 text-sm font-medium text-white transition hover:bg-gray-700"
                    >
                        Atvērt
                    </a>

                    <div className="relative">
                        <button
                            type="button"
                            onClick={copyShareLink}
                            className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-700 transition hover:bg-gray-100"
                            aria-label="Kopēt dalīšanās saiti"
                        >
                            ↗
                        </button>

                        {copied && (
                            <div
                                className="absolute bottom-full left-1/2 mb-2 -translate-x-1/2 whitespace-nowrap rounded-md bg-gray-900 px-2 py-1 text-xs font-medium text-white shadow-md">
                                Nokopēts
                            </div>
                        )}
                    </div>

                    {type === "owned" && (
                        <button
                            type="button"
                            onClick={handleDelete}
                            className="rounded-lg border border-red-300 bg-red-50 px-3 py-2 text-sm font-medium text-red-700 transition hover:bg-red-100"
                        >
                            Dzēst
                        </button>
                    )}
                </div>
            </div>
        </article>
    );
}

function formatFileSize(bytes: number) {
    if (bytes === 0) {
        return "0 B";
    }

    const units = ["B", "KB", "MB", "GB"];
    const index = Math.floor(Math.log(bytes) / Math.log(1024));

    return `${(bytes / Math.pow(1024, index)).toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}
