import gsap from 'gsap';
import { CustomEase } from 'gsap/CustomEase';
import { Observer } from 'gsap/Observer';
import { ScrollToPlugin } from 'gsap/ScrollToPlugin';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger, ScrollToPlugin, Observer, CustomEase);

CustomEase.create('opteaOut', '0.25, 0.46, 0.45, 0.94');
CustomEase.create('opteaIn', '0.55, 0.00, 1.00, 0.45');
CustomEase.create('opteaInOut', '0.76, 0.00, 0.24, 1.00');
CustomEase.create('opteaSnap', '0.34, 1.56, 0.64, 1.00');

export { gsap, Observer, ScrollTrigger };
