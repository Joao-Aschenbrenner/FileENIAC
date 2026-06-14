from PIL import Image
import io

# Read icon.ico
with open('C:/Users/USUARIO/Documents/GitWrkspc/eniac-workspace/apps/desktop/src-tauri/icons/icon.ico', 'rb') as f:
    ico_data = f.read()

ico = Image.open(io.BytesIO(ico_data))
n = 0
while True:
    try:
        ico.seek(n)
        print(f'Frame {n}: {ico.size} mode={ico.mode}')
        n += 1
    except EOFError:
        break

# Check icon.png
png = Image.open('C:/Users/USUARIO/Documents/GitWrkspc/eniac-workspace/apps/desktop/src-tauri/icons/icon.png')
print(f'icon.png: {png.size} mode={png.mode}')
# Check center pixel
cx, cy = 256, 256  # center of 512x512
pixel = png.getpixel((cx, cy))
print(f'Center pixel (should be white/light): {pixel}')
# Top-left
pixel2 = png.getpixel((10, 10))
print(f'Corner pixel (should be blue): {pixel2}')
