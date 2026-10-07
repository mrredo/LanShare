
import type {Route} from "../../.react-router/types/app/route/+types/file.delete";

    export async function clientAction({
                                           params,
        request
                                       }: Route.ClientActionArgs) {
        if(request.method !== "DELETE") return null;
    
        const response = await fetch(`/api/files/${params.id}`, {
            method: "DELETE",
            credentials: "include",
        });
    
        return await response.json();
    }