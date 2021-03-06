// SPDX-License-Identifier: Apache-2.0

package gowasm2csharp

import (
	"os"
	"path/filepath"
	"text/template"
)

type wasmData struct {
	Offset int
	Data   []byte
}

func writeMemInitData(dir string, data []wasmData) error {
	f, err := os.Create(filepath.Join(dir, "MemInitData"))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, d := range data {
		if _, err := f.WriteAt(d.Data, int64(d.Offset)); err != nil {
			return err
		}
	}
	return nil
}

func writeMemCS(dir string, namespace string, initPageNum int) error {
	f, err := os.Create(filepath.Join(dir, "Mem.cs"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := memTmpl.Execute(f, struct {
		Namespace   string
		InitPageNum int
		Data        []wasmData
	}{
		Namespace:   namespace,
		InitPageNum: initPageNum,
	}); err != nil {
		return err
	}
	return nil
}

var memTmpl = template.Must(template.New("Mem.cs").Parse(`// Code generated by go2dotnet. DO NOT EDIT.

using System;
using System.IO;
using System.Reflection;
using System.Text;

namespace {{.Namespace}}
{
    sealed class Mem
    {
        const int PageSize = 64 * 1024;

        private static void ReadFull(byte[] dst, Stream stream)
        {
            int offset = 0;
            int read = 0;
            while ((read = stream.Read(dst, offset, dst.Length - offset)) > 0)
            {
                offset += read;
            }
        }

        public Mem()
        {
            this.bytes = new byte[{{.InitPageNum}} * PageSize];

            Assembly asm = Assembly.GetExecutingAssembly();
            Stream stream = asm.GetManifestResourceStream("{{.Namespace}}.MemInitData");
            if (stream == null)
            {
                Console.Error.WriteLine("MemInitData must be embedded but not found.");
                Console.Error.WriteLine("Please add these lines under the <Project> in the csproj file:");
                Console.Error.WriteLine("");
                Console.Error.WriteLine("  <ItemGroup>");
                Console.Error.WriteLine("    <EmbeddedResource Include=\"autogen\\MemInitData\">");
                Console.Error.WriteLine("      <LogicalName>{{.Namespace}}.MemInitData</LogicalName>");
                Console.Error.WriteLine("    </EmbeddedResource>");
                Console.Error.WriteLine("  </ItemGroup>");
                Environment.Exit(1);
                return;
            }
            ReadFull(this.bytes, stream);
        }

        internal int Size
        {
            get
            {
                return this.bytes.Length / PageSize;
            }
        }

        internal int Grow(int delta)
        {
            var prevSize = this.Size;
            Array.Resize(ref this.bytes, (prevSize + delta) * PageSize);
            return prevSize;
        }

        internal sbyte LoadInt8(int addr)
        {
            return (sbyte)this.bytes[addr];
        }

        internal byte LoadUint8(int addr)
        {
            return this.bytes[addr];
        }

        internal short LoadInt16(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(short*)ptr;
                }
            }
        }

        internal ushort LoadUint16(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(ushort*)ptr;
                }
            }
        }

        internal int LoadInt32(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(int*)ptr;
                }
            }
        }

        internal uint LoadUint32(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(uint*)ptr;
                }
            }
        }

        internal long LoadInt64(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(long*)ptr;
                }
            }
        }

        internal float LoadFloat32(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(float*)ptr;
                }
            }
        }

        internal double LoadFloat64(int addr)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    return *(double*)ptr;
                }
            }
        }

        internal void StoreInt8(int addr, sbyte val)
        {
            this.bytes[addr] = (byte)val;
        }

        internal void StoreInt16(int addr, short val)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    *(short*)ptr = val;
                }
            }
        }

        internal void StoreInt32(int addr, int val)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    *(int*)ptr = val;
                }
            }
        }

        internal void StoreInt64(int addr, long val)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    *(long*)ptr = val;
                }
            }
        }

        internal void StoreFloat32(int addr, float val)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    *(float*)ptr = val;
                }
            }
        }

        internal void StoreFloat64(int addr, double val)
        {
            unsafe
            {
                fixed (byte* ptr = &this.bytes[addr])
                {
                    *(double*)ptr = val;
                }
            }
        }

        internal void StoreBytes(int addr, byte[] bytes)
        {
            Array.Copy(bytes, 0, this.bytes, addr, bytes.Length);
        }

        internal ArraySegment<byte> LoadSlice(int addr)
        {
            var array = this.LoadInt64(addr);
            var len = this.LoadInt64(addr + 8);
            return new ArraySegment<byte>(this.bytes, (int)array, (int)len);
        }

        internal ArraySegment<byte> LoadSliceDirectly(long array, int len)
        {
            return new ArraySegment<byte>(this.bytes, (int)array, len);
        }

        internal string LoadString(int addr)
        {
            var saddr = this.LoadInt64(addr);
            var len = this.LoadInt64(addr + 8);
            return Encoding.UTF8.GetString(this.bytes, (int)saddr, (int)len);
        }

        private byte[] bytes;
    }
}
`))
