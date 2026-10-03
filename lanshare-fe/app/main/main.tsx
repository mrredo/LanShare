import { useState } from "react";
import type { file } from "~/types/file";
import LanShareTitle from "~/components/title";
import Background from "~/components/Background";

export function meta() {
  return [
    { title: "LanShare" },
    { name: "description", content: "LanShare failu dalīšanās" },
  ];
}

const myFilesC: file[] = [
  {
    id: "my-001",
    filename: "project.zip",
    storage_path: "/uploads/project.zip",
    size: 2457600,
    owner_cookie: "user-cookie",
    uploaded_at: "2026-10-03T09:15:00Z",
    expires_at: "2026-10-04T09:15:00Z",
  },
  {
    id: "my-002",
    filename: "notes.txt",
    storage_path: "/uploads/notes.txt",
    size: 4820,
    owner_cookie: "user-cookie",
    text: "Important project notes",
    uploaded_at: "2026-10-02T14:30:00Z",
    expires_at: "2026-10-05T14:30:00Z",
  },
  {
    id: "my-003",
    filename: "presentation.pdf",
    storage_path: "/uploads/presentation.pdf",
    size: 1843200,
    owner_cookie: "user-cookie",
    uploaded_at: "2026-10-01T11:20:00Z",
    expires_at: "2026-10-08T11:20:00Z",
  },
];

const publicFilesC: file[] = [
  {
    id: "public-001",
    filename: "README.md",
    storage_path: "/uploads/README.md",
    size: 3200,
    owner_cookie: "other-user",
    text: "# LanShare\nWelcome to LanShare!",
    uploaded_at: "2026-10-03T08:00:00Z",
    expires_at: "2026-10-04T08:00:00Z",
  },
  {
    id: "public-002",
    filename: "image.png",
    storage_path: "/uploads/image.png",
    size: 5242880,
    owner_cookie: "other-user",
    uploaded_at: "2026-10-02T16:45:00Z",
    expires_at: "2026-10-09T16:45:00Z",
  },
  {
    id: "public-003",
    filename: "website.url",
    storage_path: "/uploads/website.url",
    size: 64,
    owner_cookie: "another-user",
    url: "https://example.com",
    uploaded_at: "2026-10-01T10:10:00Z",
    expires_at: "2026-10-06T10:10:00Z",
  },
];

export default function Main() {
  const [myFiles] = useState<file[]>(myFilesC);
  const [publicFiles] = useState<file[]>(publicFilesC);

  return (
      <div className="relative min-h-screen w-full overflow-hidden px-5 py-12">
        <Background />

        <main className="relative z-10 mx-auto max-w-7xl">
          <header className="mb-12 text-center">
            <LanShareTitle />

            <p className="mt-4 text-lg text-gray-600">
              Ērti kopīgo un pārvaldi failus lokālajā tīklā.
            </p>
          </header>

          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            <FileListSection
                files={myFiles}
                title="Mani faili"
                description="Tavi augšupielādētie faili"
                type="owned"
            />

            <FileListSection
                files={publicFiles}
                title="Publiskie faili"
                description="Faili, ko kopīgojuši citi lietotāji"
                type="public"
            />
          </div>
        </main>
      </div>
  );
}

export function FileListSection({
                                  files,
                                  title,
                                  description,
                                  type,
                                }: {
  files: file[];
  title: string;
  description: string;
  type: "owned" | "public";
}) {
  return (
      <section className="rounded-2xl border border-gray-400/70 bg-gray-200/80 p-5 shadow-lg backdrop-blur-sm">
        <div className="mb-5 flex items-start justify-between gap-4">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">{title}</h2>
            <p className="mt-1 text-sm text-gray-600">{description}</p>
          </div>

          <span className="rounded-full bg-gray-300 px-3 py-1 text-sm font-medium text-gray-700">
          {files.length}
        </span>
        </div>

        <div className="grid grid-cols-1 gap-3">
          {files.map((file) => (
              <FileCard key={file.id} file={file} type={type} />
          ))}
        </div>
      </section>
  );
}

export function FileCard({
                           file,
                           type,
                         }: {
  file: file;
  type: "owned" | "public";
}) {
  let fileName = file.filename;

  if (!fileName) {
    fileName = file.text?.slice(0, 50);
  }

  if (!fileName) {
    fileName = file.url;
  }

  return (
      <article className="group rounded-xl border border-gray-300/80 bg-gray-100/90 p-4 shadow-sm backdrop-blur-sm transition-all duration-200 hover:-translate-y-0.5 hover:border-gray-400 hover:bg-white hover:shadow-md">
        <div className="flex items-center justify-between gap-4">
          <a href={`/files/${file.id}`} className="min-w-0 flex-1">
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

          <div className="flex shrink-0 items-center gap-2 opacity-0 transition-opacity duration-200 group-hover:opacity-100">
            <a
                href={`/files/${file.id}`}
                className="rounded-lg bg-gray-800 px-3 py-2 text-sm font-medium text-white transition hover:bg-gray-700"
            >
              Atvērt
            </a>

            {type === "owned" && (
                <button
                    type="button"
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