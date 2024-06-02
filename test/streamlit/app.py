import os
import streamlit as st

def main():
    st.title("Simple Streamlit App")
    user_input = st.text_input("Enter some text")
    
    if st.button("Submit"):
        st.write(f"You entered: {user_input}: {os.getpid()}")
        for key, value in os.environ.items():
            st.write(f'{key}: {value}')

if __name__ == "__main__":
    main()

