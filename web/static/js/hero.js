import * as THREE from 'three';
import { GLTFLoader } from 'three/addons/loaders/GLTFLoader.js';
import { OrbitControls } from 'three/addons/controls/OrbitControls.js';

const portrait = document.getElementById('char-portrait');
if (!portrait) throw new Error('no portrait');

const canvas  = document.getElementById('rpg-canvas');
const loading = document.getElementById('canvas-loading');
const label   = document.getElementById('canvas-class-label');
const hint    = label ? label.nextElementSibling : null;
const charClass = (portrait.dataset.class || '').toLowerCase();

const classToModel = {
    guardian: 'Knight', berserker: 'Barbarian', paladin: 'Knight',
    rogue: 'Rogue_Hooded', sage: 'Mage', knight: 'Knight',
    battlemage: 'Mage', warlord: 'Barbarian', wanderer: 'Ranger'
};

const ALL_MODELS  = ['Mage', 'Knight', 'Barbarian', 'Ranger', 'Rogue_Hooded'];
const MODEL_LABEL = { Mage:'SAGE', Knight:'GUARDIAN', Barbarian:'BERSERKER', Ranger:'WANDERER', Rogue_Hooded:'ROGUE' };
const startModel  = classToModel[charClass] || 'Mage';
let currentIdx    = ALL_MODELS.indexOf(startModel);
if (currentIdx < 0) currentIdx = 0;

if (label) label.textContent = MODEL_LABEL[startModel] || 'MAGE';
if (hint)  hint.textContent  = 'CLICK TO SWITCH ▶';

const W = 260, H = 340;
const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
renderer.setSize(W, H);
renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
renderer.outputColorSpace = THREE.SRGBColorSpace;

const scene = new THREE.Scene();
scene.background = new THREE.Color(0x050010);

const camera = new THREE.PerspectiveCamera(45, W / H, 0.1, 100);
camera.position.set(0, 1.5, 3.8);

const controls = new OrbitControls(camera, canvas);
controls.target.set(0, 1.0, 0);
controls.enablePan = false;
controls.enableZoom = false;
controls.minPolarAngle = Math.PI * 0.25;
controls.maxPolarAngle = Math.PI * 0.75;
controls.autoRotate = true;
controls.autoRotateSpeed = 2.0;
controls.enableDamping = true;
controls.dampingFactor = 0.08;
controls.update();

scene.add(new THREE.AmbientLight(0x5522aa, 1.2));
const key = new THREE.DirectionalLight(0xFFD700, 2.2);
key.position.set(2, 5, 3);
scene.add(key);
const fill = new THREE.DirectionalLight(0x00E5FF, 0.7);
fill.position.set(-2, 1, -2);
scene.add(fill);

let model = null;
const loader = new GLTFLoader();

function normalizeModel(m) {
    const box = new THREE.Box3().setFromObject(m);
    const size = new THREE.Vector3();
    box.getSize(size);
    const s = 2.0 / size.y;
    m.scale.setScalar(s);
    box.setFromObject(m);
    const center = new THREE.Vector3();
    box.getCenter(center);
    m.position.set(-center.x, -box.min.y, -center.z);
}

function loadModel(name) {
    if (model) { scene.remove(model); model = null; }
    if (loading) loading.style.display = 'flex';
    loader.load('/static/assets/3d/characters/' + name + '.glb', function(gltf) {
        model = gltf.scene;
        normalizeModel(model);
        scene.add(model);
        if (loading) loading.style.display = 'none';
        if (label) label.textContent = MODEL_LABEL[name] || name;
    }, undefined, function() {
        if (loading) loading.style.display = 'none';
    });
}

loadModel(ALL_MODELS[currentIdx]);

canvas.style.cursor = 'grab';

let dragMoved = false;
canvas.addEventListener('mousedown', function() { dragMoved = false; });
canvas.addEventListener('mousemove', function() { dragMoved = true; });
canvas.addEventListener('mouseup', function() {
    if (!dragMoved) {
        currentIdx = (currentIdx + 1) % ALL_MODELS.length;
        loadModel(ALL_MODELS[currentIdx]);
    }
});
canvas.addEventListener('mousedown', function() { canvas.style.cursor = 'grabbing'; });
canvas.addEventListener('mouseup',   function() { canvas.style.cursor = 'grab'; });

let animFrameId;
(function animate() {
    animFrameId = requestAnimationFrame(animate);
    controls.update();
    renderer.render(scene, camera);
}());

document.addEventListener('htmx:beforeSwap', function() {
    cancelAnimationFrame(animFrameId);
    controls.dispose();
    renderer.dispose();
}, { once: true });
