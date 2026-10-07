import { useEffect, useRef } from "react";
import gsap from "gsap";

export default function LanShareTitle() {
    const titleRef = useRef<HTMLHeadingElement>(null);

    useEffect(() => {
        const title = titleRef.current;

        if (!title) {
            return;
        }

        const text = title.textContent ?? "";
        title.innerHTML = "";

        [...text].forEach((char) => {
            const span = document.createElement("span");
            span.textContent = char === " " ? "\u00A0" : char;
            span.style.display = "inline-block";
            title.appendChild(span);
        });

        const chars = title.querySelectorAll("span");

        gsap.set(chars, {
            y: 100,
            opacity: 0,
            rotateX: -90,
            filter: "blur(10px)",
            transformOrigin: "50% 100%",
        });

        const animation = gsap.to(chars, {
            y: 0,
            opacity: 1,
            rotateX: 0,
            filter: "blur(0px)",
            duration: 0.9,
            stagger: {
                each: 0.035,
                from: "center",
            },
            ease: "back.out(1.7)",
        });

        return () => {
            animation.kill();
        };
    }, []);

    return (
        <h1
            ref={titleRef}
            className="w-full text-center text-5xl font-black tracking-tight text-gray-900 md:text-6xl"
            style={{ perspective: "800px" }}
        >
            LanShare failu dalīšanās
        </h1>
    );
}