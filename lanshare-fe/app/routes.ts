import {type RouteConfig, index, route} from "@react-router/dev/routes";

export default [
    index("route/main.tsx"),
    route("/upload", "route/fileUpload.tsx"),
    route("/files/:id", "route/fileView.tsx"),

    route("/files/delete/:id", "route/file.delete.tsx")


] satisfies RouteConfig;
