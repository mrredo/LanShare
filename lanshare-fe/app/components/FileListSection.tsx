import type { file } from "~/types/file";
import FileCard from "~/components/FileCard";

export default function FileListSection({
                                            files,
                                            title,
                                            description,
                                            type,
                                            route,
                                            error,
                                            showAddButton = false,
                                        }: {
    files: file[];
    title: string;
    description: string;
    type: "owned" | "public";
    route: "/files" | "/share";
    error?: string;
    showAddButton?: boolean;
}) {
    return (
        <section className="rounded-2xl border border-gray-400/70 bg-gray-200/80 p-5 shadow-lg backdrop-blur-sm">
            <div className="mb-5 flex items-start justify-between gap-4">
                <div>
                    <h2 className="text-2xl font-bold text-gray-900">{title}</h2>
                    <p className="mt-1 text-sm text-gray-600">{description}</p>
                </div>

                <div className="flex items-center gap-2">
                    <span className="rounded-full bg-gray-300 px-3 py-1 text-sm font-medium text-gray-700">
                        {files.length}
                    </span>

                    {showAddButton && (
                        <a
                            href="/upload"
                            className="flex h-8 w-8 items-center justify-center rounded-full bg-gray-800 text-xl font-medium leading-none text-white transition hover:bg-gray-700"
                            title="Pievienot failu"
                            aria-label="Pievienot failu"
                        >
                            +
                        </a>
                    )}
                </div>
            </div>

            {error && (
                <div className="flex h-full w-full items-center justify-center">
                    <p className="text-red-500">{error}</p>
                </div>
            )}

            {!error && files.length === 0 && (
                <div className="flex  w-full items-center justify-center">
                    <p className="font-bold">Nav failu...</p>
                </div>
            )}

            {!error && files.length > 0 && (
                <div className="grid grid-cols-1 gap-3">
                    {files.map((file) => (
                        <FileCard
                            route={route}
                            key={file.id}
                            file={file}
                            type={type}
                        />
                    ))}
                </div>
            )}
        </section>
    );
}