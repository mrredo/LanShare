import { useEffect, useRef } from "react";
import gsap from "gsap";

export default function Background() {
    const backgroundRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const context = gsap.context(() => {
            const elements = gsap.utils.toArray<HTMLElement>(".background-element");

            elements.forEach((element, index) => {
                gsap.to(element, {
                    x: gsap.utils.random(-180, 180),
                    y: gsap.utils.random(-140, 140),
                    rotation: gsap.utils.random(-180, 180),
                    scale: gsap.utils.random(0.8, 1.3),
                    duration: gsap.utils.random(12, 20),
                    repeat: -1,
                    yoyo: true,
                    ease: "sine.inOut",
                    delay: index * -2,
                });
            });
        }, backgroundRef);

        return () => context.revert();
    }, []);

    return (
        <div
            ref={backgroundRef}
            className="pointer-events-none fixed inset-0 -z-10 overflow-hidden bg-gray-300"
        >
            <div className="background-element absolute -left-32 -top-32 h-96 w-96 rounded-full bg-gray-400/40 blur-3xl" />

            <div className="background-element absolute right-[-10rem] top-[15%] h-[30rem] w-[30rem] rounded-full bg-gray-400/30 blur-3xl" />

            <div className="background-element absolute bottom-[-12rem] left-[15%] h-[32rem] w-[32rem] rounded-full bg-gray-400/40 blur-3xl" />

            <div className="background-element absolute bottom-[10%] right-[20%] h-72 w-72 rounded-full bg-gray-400/25 blur-3xl" />

            <div className="background-element absolute left-[40%] top-[20%] h-48 w-48 rounded-full bg-white/20 blur-3xl" />
        </div>
    );
}