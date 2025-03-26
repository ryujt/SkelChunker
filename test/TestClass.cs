using System;
using System.Collections.Generic;
using System.Linq;

namespace TestNamespace
{
    // 테스트 클래스 주석
    public class TestClass
    {
        // 프로퍼티
        public int Id { get; set; }
        public string Name { get; private set; }
        
        // 생성자
        public TestClass(string name)
        {
            Name = name;
            Id = 0;
        }
        
        // 메서드
        public void PrintInfo()
        {
            Console.WriteLine($"Id: {Id}, Name: {Name}");
        }
        
        // 중첩 클래스
        public class NestedClass
        {
            public string Description { get; set; }
            
            public void ShowDescription()
            {
                Console.WriteLine(Description);
            }
        }
    }
    
    // 인터페이스
    public interface ITestInterface
    {
        void TestMethod();
        string TestProperty { get; }
    }
    
    // 구조체
    public struct TestStruct
    {
        public int X1;
        public int Y2;
        
        public TestStruct(int x, int y)
        {
            X1 = x;
            Y2 = y;
        }
        
        public double GetDistance()
        {
            return Math.Sqrt(X1 * X1 + Y2 * Y2);
        }
    }
    
    // 레코드
    public record TestRecord(string Title, string Description);
} 