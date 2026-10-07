import { useState } from "react";
import type { file } from "~/types/file";
import LanShareTitle from "~/components/title";
import Background from "~/components/Background";
import FileListSection from "~/components/FileListSection";
import {useLoaderData} from "react-router";

export function meta() {
  return [
    { title: "LanShare" },
    { name: "description", content: "LanShare failu dalīšanās" },
  ];
}
export async function clientLoader() {
  const res = await fetch("/api/files", {
    method: "GET",
    credentials: "include",
  })
  return await res.json();
}
type LoaderData = {
  public_files_error?: string;
  user_files_error?: string;
  public_files?: file[];
  user_files?: file[];
}
export function HydrateFallback() {
  return <p>Meklējam failus</p>;
}
export default function Main() {
  const loaderData = useLoaderData<LoaderData>();
  const publicFiles = loaderData.public_files ?? [];
  const myFiles = loaderData.user_files ?? [];
  const publicFileError = loaderData.public_files_error;
  const userFileError = loaderData.user_files_error;


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
                route={"/files"}
                files={myFiles}
                error={userFileError}
                showAddButton={true}
                title="Mani faili"
                description="Tavi augšupielādētie faili"
                type="owned"
            />

            <FileListSection
                route={"/share"}
                files={publicFiles}
                error={publicFileError}
                title="Publiskie faili"
                description="Faili, ko kopīgojuši citi lietotāji"
                type="public"
            />
          </div>
        </main>
      </div>
  );
}
